package main

import (
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const FINAL_DIR = "../dataJoin/"

func main() {
	cntImg := flag.Int("cntImg", 10000000, "Cantidad de imagenes a copiar")
	resize := flag.Bool("resize", false, "Si es que se quieren rescalar la imagen a 950x544")
	flag.Parse()
	cnt := 0

	checkDirectorio(FINAL_DIR)
	l := listDirectoryRecursive("../dataset_1088x612")
	for e := l.Front(); e != nil; e = e.Next() {
		if cnt += 1; cnt <= *cntImg {
			newName := getNewName(e.Value.(string))
			ext_aux := strings.Split(newName, ".")
			ext := strings.ToLower(ext_aux[len(ext_aux)-1])

			if *resize {
				dir, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				vipsExec(dir+"/"+e.Value.(string), dir+"/"+FINAL_DIR+newName)
				//resizeLabels(strings.Replace(e.Value.(string), "."+ext, ".txt", 1))
			} else {
				err := CopyFile(e.Value.(string), FINAL_DIR+newName)
				if err != nil {
					log.Fatal(err)
				}
			}
			err := CopyFile(strings.Replace(e.Value.(string), "."+ext, ".txt", 1),
				FINAL_DIR+strings.Replace(newName, "."+ext, ".txt", 1))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func resizeLabels(src string) {

	kitti_line := "{1} 0.00 0 0.00 {2} {3} {4} {5} 0.00 0.00 0.00 0.00 0.00 0.00 0.00"
	file, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		labels_split := strings.Split(fileScanner.Text(), " ")
		label := " "
		if labels_split[0] == "0" {
			label = "sinCasco"
		} else if labels_split[0] == "1" {
			label = "conCasco"
		} else if labels_split[0] == "2" {
			label = "proteccionAuditiva"
		} else if labels_split[0] == "3" {
			label = "mascaraSoldador"
		} else if labels_split[0] == "4" {
			label = "pechoDesnudo"
		} else if labels_split[0] == "5" {
			label = "chalecoReflectante"
		} else if labels_split[0] == "6" {
			label = "persona"
		}

		a, _ := strconv.ParseFloat(labels_split[1], 32)
		b, _ := strconv.ParseFloat(labels_split[2], 32)
		c, _ := strconv.ParseFloat(labels_split[3], 32)
		d, _ := strconv.ParseFloat(labels_split[4], 32)

		xmin := (a - 0.5*c) * 544
		ymin := (b - 0.5*d) * 950
		xmax := (a + 0.5*c) * 544
		ymax := (b + 0.5*d) * 950

		if xmin < 0 {
			xmin = 0
		}
		if ymin < 0 {
			ymin = 0
		}
		if xmax > 544 {
			xmax = 544
		}
		if ymax > 950 {
			ymax = 950
		}
		xmin_ := strconv.FormatFloat(xmin, 'E', -1, 32)
		ymin_ := strconv.FormatFloat(ymin, 'E', -1, 32)
		xmax_ := strconv.FormatFloat(xmax, 'E', -1, 32)
		ymax_ := strconv.FormatFloat(ymax, 'E', -1, 32)

		kitti_lineWrite := strings.Replace(kitti_line, "{1}", label, 1)
		kitti_lineWrite = strings.Replace(kitti_lineWrite, "{2}", xmin_, 1)
		kitti_lineWrite = strings.Replace(kitti_lineWrite, "{3}", ymin_, 1)
		kitti_lineWrite = strings.Replace(kitti_lineWrite, "{4}", xmax_, 1)
		kitti_lineWrite = strings.Replace(kitti_lineWrite, "{5}", ymax_, 1)
		fmt.Println(kitti_lineWrite)
	}
	file.Close()
}

func listDirectoryRecursive(src string) (l_images *list.List) {
	l_img := list.New()
	archivos, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatal(err)
	}
	for _, archivo := range archivos {
		if archivo.IsDir() {
			l_img.PushBackList(listDirectoryRecursive(src + "/" + archivo.Name()))
		} else {
			split_r := strings.Split(archivo.Name(), ".")
			extension := strings.ToLower(split_r[len(split_r)-1])
			if extension == "png" || extension == "jpg" || extension == "jpeg" {
				l_img.PushBack(src + "/" + archivo.Name())
			}
		}
	}
	return l_img
}

func vipsExec(src string, dst string) {
	var args = []string{
		"-s", "950x544!",
		"-o", dst,
		src,
	}
	if runtime.GOOS == "linux" {
		err := exec.Command("vipsthumbnail", args...).Run()
		if err != nil {
			log.Fatal("Error")
		}
	}
}

func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func checkDirectorio(directorio string) {
	if _, err := os.Stat(directorio); !os.IsNotExist(err) {
		err = os.RemoveAll(directorio)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := os.Mkdir(directorio, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func getNewName(src string) (output string) {
	splitLinea := strings.Split(src, "/")
	newName := splitLinea[len(splitLinea)-3] + "_" + splitLinea[len(splitLinea)-2] + "_" + splitLinea[len(splitLinea)-1]
	return newName
}
