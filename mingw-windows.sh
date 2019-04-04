#!/bin/sh

#  dependencies

echo "Installing dependencies via Homebrew (http://brew.sh)"

# ruby -e "$(curl -fsSL https://raw.github.com/Homebrew/homebrew/go/install)"
brew update
brew tap homebrew/versions

brew install gcc5
brew install wget

#  mingw

VERSION="4.0.4"
BINUTILS_VERSION="2.25.1"
GCC_MAJOR="5"
GCC_VERSION="5.2.0"
HOMEBREW_PREFIX=`brew --prefix`
PREFIX="${HOMEBREW_PREFIX}/Cellar/mingw-w64/${VERSION}"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

mkdir -p source
mkdir -p $PREFIX

echo "Downloading binutils\n"

cd ./source
wget -O binutils-${BINUTILS_VERSION}.tar.bz2 http://ftpmirror.gnu.org/binutils/binutils-${BINUTILS_VERSION}.tar.bz2
tar xjf binutils-${BINUTILS_VERSION}.tar.bz2

echo "Building binutils\n"
echo "1/2 32-bit\n"

cd binutils-${BINUTILS_VERSION}
mkdir build
cd build

CC=gcc-${GCC_MAJOR} CXX=g++-${GCC_MAJOR} CPP=cpp-${GCC_MAJOR} LD=gcc-${GCC_MAJOR} ../configure --target=i686-w64-mingw32 --disable-werror --disable-multilib --prefix=$PREFIX --with-sysroot=$PREFIX
make -j4
make install-strip

echo "2/2 64-bit\n"
cd ..
rm -rf build
mkdir build
cd build

CC=gcc-${GCC_MAJOR} CXX=g++-${GCC_MAJOR} CPP=cpp-${GCC_MAJOR} LD=gcc-${GCC_MAJOR} ../configure --target=x86_64-w64-mingw32 --disable-werror --disable-multilib --prefix=$PREFIX --with-sysroot=$PREFIX --enable-64-bit-bfd
make -j4
make install-strip

cd ..
cd ..

echo "Downloading mingw-w64\n"

wget -O mingw-w64-v${VERSION}.tar.bz2 "https://downloads.sourceforge.net/project/mingw-w64/mingw-w64/mingw-w64-release/mingw-w64-v${VERSION}.tar.bz2"
tar xjf mingw-w64-v${VERSION}.tar.bz2

echo "Building mingw-headers\n"

echo "1/2 32-bit\n"

cd mingw-w64-v${VERSION}
mkdir build-headers
cd build-headers

../mingw-w64-headers/configure --host=i686-w64-mingw32 --prefix=$PREFIX/i686-w64-mingw32
make -j4
make install-strip

cd $PREFIX/i686-w64-mingw32
ln -s lib lib32
cd $DIR/source/mingw-w64-v${VERSION}

echo "2/2 64-bit\n"
rm -rf build-headers
mkdir build-headers
cd build-headers

../mingw-w64-headers/configure --host=x86_64-w64-mingw32 --prefix=$PREFIX/x86_64-w64-mingw32
make -j4
make install-strip


cd $PREFIX/x86_64-w64-mingw32
ln -s lib lib64
cd $DIR/source/

echo "Downloading gcc\n"

wget -O gcc-${GCC_VERSION}.tar.bz2 http://ftpmirror.gnu.org/gcc/gcc-${GCC_VERSION}/gcc-${GCC_VERSION}.tar.bz2
tar xjf gcc-${GCC_VERSION}.tar.bz2

echo "Building core gcc\n"

echo "1/2 32-bit\n"


cd $PREFIX
ln -s i686-w64-mingw32 mingw

cd $DIR/source/gcc-${GCC_VERSION}
mkdir build32
cd build32

CC=gcc-${GCC_MAJOR} CXX=g++-${GCC_MAJOR} CPP=cpp-${GCC_MAJOR} LD=gcc-${GCC_MAJOR} PATH=${HOMEBREW_PREFIX}/mingw/bin/:$PATH ../configure --target=i686-w64-mingw32 --disable-multilib --enable-languages=c,c++,objc,obj-c++ --with-gmp=${HOMEBREW_PREFIX}/opt/gmp/ --with-mpfr=${HOMEBREW_PREFIX}/opt/mpfr/ --with-mpc=${HOMEBREW_PREFIX}/opt/libmpc/ --with-isl=${HOMEBREW_PREFIX}/opt/isl014/ --with-system-zlib --enable-version-specific-runtime-libs --enable-libstdcxx-time=yes --enable-stage1-checking --enable-checking=release --enable-lto --enable-threads=win32 --disable-sjlj-exceptions --prefix=$PREFIX --with-sysroot=$PREFIX

PATH=${PREFIX}/bin/:$PATH make all-gcc -j4
PATH=${PREFIX}/bin/:$PATH make install-gcc

echo "2/2 64-bit\n"

cd $PREFIX
rm mingw
ln -s x86_64-w64-mingw32 mingw

cd $DIR/source/gcc-${GCC_VERSION}
mkdir build64
cd build64

CC=gcc-${GCC_MAJOR} CXX=g++-${GCC_MAJOR} CPP=cpp-${GCC_MAJOR} LD=gcc-${GCC_MAJOR} PATH=${HOMEBREW_PREFIX}/mingw/bin/:$PATH ../configure --target=x86_64-w64-mingw32 --disable-multilib --enable-languages=c,c++,objc,obj-c++ --with-gmp=${HOMEBREW_PREFIX}/opt/gmp/ --with-mpfr=${HOMEBREW_PREFIX}/opt/mpfr/ --with-mpc=${HOMEBREW_PREFIX}/opt/libmpc/ --with-isl=${HOMEBREW_PREFIX}/opt/isl014/ --with-system-zlib --enable-version-specific-runtime-libs --enable-libstdcxx-time=yes --enable-stage1-checking --enable-checking=release --enable-lto --enable-threads=win32 --prefix=$PREFIX --with-sysroot=$PREFIX

PATH=${PREFIX}/bin/:$PATH make all-gcc -j4
PATH=${PREFIX}/bin/:$PATH make install-gcc

echo "Building mingw runtime\n"

cd $PREFIX
rm mingw
ln -s i686-w64-mingw32 mingw

cd $DIR/source/mingw-w64-v${VERSION}
mkdir build-crt
cd build-crt

echo "1/2 32-Bit\n"

PATH=${PREFIX}/bin/:$PATH ../mingw-w64-crt/configure --host=i686-w64-mingw32 --prefix=$PREFIX/i686-w64-mingw32 --with-sysroot=$PREFIX

PATH=${PREFIX}/bin/:$PATH make
PATH=${PREFIX}/bin/:$PATH make install-strip

echo "2/2 64-Bit\n"

cd $PREFIX
rm mingw
ln -s x86_64-w64-mingw32 mingw

cd $DIR/source/mingw-w64-v${VERSION}
rm -rf build-crt
mkdir build-crt
cd build-crt

PATH=${PREFIX}/bin/:$PATH ../mingw-w64-crt/configure --host=x86_64-w64-mingw32 --prefix=$PREFIX/x86_64-w64-mingw32 --with-sysroot=$PREFIX

PATH=${PREFIX}/bin/:$PATH make
PATH=${PREFIX}/bin/:$PATH make install-strip

echo "Building all gcc\n"

echo "1/2 32-Bit\n"

cd $PREFIX
rm mingw
ln -s i686-w64-mingw32 mingw

cd $DIR/source/gcc-${GCC_VERSION}/build32

PATH=${PREFIX}/bin/:$PATH make
PATH=${PREFIX}/bin/:$PATH make install-strip

echo "2/2 64-Bit\n"

cd $PREFIX
rm mingw
ln -s x86_64-w64-mingw32 mingw

cd $DIR/source/gcc-${GCC_VERSION}/build64

PATH=${PREFIX}/bin/:$PATH make
PATH=${PREFIX}/bin/:$PATH make install-strip

echo "Linking libgcc\n"

cd $PREFIX/i686-w64-mingw32/lib
ln -s ../../lib/gcc/i686-w64-mingw32/lib/libgcc_s.a ./

cd $PREFIX/x86_64-w64-mingw32/lib
ln -s ../../lib/gcc/x86_64-w64-mingw32/lib/libgcc_s.a ./

echo "Building winpthreads\n"

cd $DIR/source/mingw-w64-v${VERSION}/mingw-w64-libraries/winpthreads

echo "1/2 32-Bit\n"

mkdir -p build
cd build

PATH=${PREFIX}/bin/:$PATH ../configure --host=i686-w64-mingw32 --prefix=$PREFIX/i686-w64-mingw32
PATH=${PREFIX}/bin/:$PATH make
PATH=${PREFIX}/bin/:$PATH make install-strip

cd ..
rm -rf build

echo "2/2 64-Bit\n"

mkdir -p build
cd build

PATH=${PREFIX}/bin/:$PATH ../configure --host=x86_64-w64-mingw32 --prefix=$PREFIX/x86_64-w64-mingw32
PATH=${PREFIX}/bin/:$PATH make
PATH=${PREFIX}/bin/:$PATH make install-strip

echo "Cleaning up\n"

cd $DIR
rm -rf source

echo "Done"