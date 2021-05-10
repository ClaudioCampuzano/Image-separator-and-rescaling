# Image-separator-and-rescaling
From the images in a directory, separate and rescale them to a defined size.

# Requirements
```
wget https://github.com/libvips/libvips/releases/download/v8.10.6/vips-8.10.6.tar.gz
tar xf vips-8.10.6.tar.gz
cd vips-8.10.6/
sudo apt-get install build-essential pkg-config glib2.0-dev libexpat1-dev
```

You’ll need the dev packages for the file format support you want. For basic jpeg and tiff support, you’ll need libtiff5-dev, libjpeg-turbo8-dev, and libgsf-1-dev. See the Dependencies section below for a full list of the things that libvips can be configured to use.
In https://libvips.github.io/libvips/install.html

```
make
sudo make install
sudo ldconfig
```

# To run

```
git clone https://github.com/ClaudioCampuzano/Image-separator-and-rescaling.git
cd Image-separator-and-rescaling
go mod tidy
go run separadoV2.go -cntImg 20 -resize
```
