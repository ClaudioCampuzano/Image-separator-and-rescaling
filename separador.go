package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"os/exec"
	"runtime"
	"strings"
	//"image/png"
	//"github.com/nfnt/resize"
	//"github.com/disintegration/imaging"
)

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

//si existe lo destruye,
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

//const FINAL_DIR = "/home/ubuntu/Downloads/virtual_dataset/dataJoin_ImageResize/"

const FINAL_DIR_IMAGES = "../dataJoin_images/"
const FINAL_DIR_LABELS = "../dataJoin_labels/"

func main() {
	/*sourceDir := flag.String("sourceDir", "/home/carpetSrc", "Directorio desde donde se estraeran las imagenes")
	finalDir := flag.String("finalDir", "/home/carpetFinal", "Directorio destino de las imagenes")
	resize := flag.String("resizeImg", "-", "Tamaño al que se quieren rescalar las imagenes 950x544, si no se indica no se escalan")
	*/
	cntImg := flag.Int("cntImg", 10000000, "Cantidad de imagenes a copiar")
	flag.Parse()
	cnt := 0

	checkDirectorio(FINAL_DIR_IMAGES)
	checkDirectorio(FINAL_DIR_LABELS)
	archivosTexto := [...]string{"../dataset_1088x612/train.virtual.txt", "../dataset_1088x612/valid.virtual.txt"}
	for _, name := range archivosTexto {
		archivo, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		defer archivo.Close()

		scanner := bufio.NewScanner(archivo)
		for scanner.Scan() {
			if cnt += 1; cnt <= *cntImg {
				linea := scanner.Text()
				splitLinea := strings.Split(linea, "/")
				newName := splitLinea[2] + "_" + splitLinea[3] + "_" + splitLinea[4]
				//Resize imagenes
				//vipsExec("/home/ubuntu/Downloads/virtual_dataset"+linea, FINAL_DIR_IMAGES+newName)
				err := CopyFile(".."+linea, FINAL_DIR_IMAGES+newName)
				if err != nil {
					log.Fatal(err)
				}
				err1 := CopyFile(".."+strings.Replace(linea, ".png", ".txt", 1),
					FINAL_DIR_LABELS+strings.Replace(newName, ".png", ".txt", 1))
				if err1 != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
