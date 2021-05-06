package main

import (
	"bufio"
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

const FINAL_DIR = "/home/ubuntu/Downloads/virtual_dataset/dataJoin_ImageResize/"

//const FINAL_DIR_IMAGES = "../dataJoin_images/"
//const FINAL_DIR_LABELS = "../dataJoin_labels/"

func main() {

	checkDirectorio(FINAL_DIR)
	archivosTexto := [...]string{"../dataset_1088x612/train.virtual.txt", "../dataset_1088x612/valid.virtual.txt"}
	for _, name := range archivosTexto {
		archivo, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		defer archivo.Close()

		scanner := bufio.NewScanner(archivo)
		for scanner.Scan() {
			linea := scanner.Text()
			splitLinea := strings.Split(linea, "/")

			newName := splitLinea[2] + "_" + splitLinea[3] + "_" + splitLinea[4]

			vipsExec("/home/ubuntu/Downloads/virtual_dataset"+linea, FINAL_DIR+newName)

			//vipsthumbnail(".."+linea, FINAL_DIR+newName)

			//resize imagenes pero lento
			/*file, err := os.Open(".." + linea)
			if err != nil {
				log.Fatal(err)
			}

			// decode jpeg into image.Image
			img, err := png.Decode(file)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()

			// resize to width 1000 using Lanczos resampling
			// and preserve aspect ratio
			m := resize.Resize(960, 544, img, resize.Lanczos3)

			out, err := os.Create(FINAL_DIR + newName)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			// write new image to file
			png.Encode(out, m)

			//resize imagenes pero lento
			src, err := imaging.Open(".." + linea)
			if err != nil {
				log.Fatalf("failed to open image: %v", err)
			}

			src_Resize := imaging.Resize(src, 960, 544, imaging.Lanczos)

			err = imaging.Save(src_Resize, FINAL_DIR+newName)
			if err != nil {
				log.Fatalf("failed to save image: %v", err)
			}

			///Copiar y pegar imagenes
			/*err := CopyFile(".."+linea, FINAL_DIR+newName)
			if err != nil {
				log.Fatal(err)
			}
			err1 := CopyFile(".."+strings.Replace(linea, ".png", ".txt", 1),
				FINAL_DIR+strings.Replace(newName, ".png", ".txt", 1))
			if err1 != nil {
				log.Fatal(err)
			}*/
		}
	}
}
