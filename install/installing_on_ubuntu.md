#   Remove go
```
sudo apt remove --auto-remove -y golang
sudo apt remove --auto-remove -y golang-go
```

# Install go=1.19
```
sudo apt update
sudo apt-get install -y \
    build-essential \
    libseccomp-dev \
    pkg-config \
    squashfs-tools \
    squashfuse \
    fuse2fs \
    fuse-overlayfs \
    fakeroot \
    cryptsetup \
    curl wget git
```
```
sudo apt install uidmap
wget -O /tmp/go1.19.linux-amd64.tar.gz https://dl.google.com/go/go1.19.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf /tmp/go1.19.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

# Install apptainer
```
git clone https://github.com/apptainer/apptainer.git
cd apptainer
#maybe not: `git checkout v1.0.3`
#maybe: `git checkout v1.1.0`
./mconfig
cd ./builddir
make
sudo make install
apptainer --version
```

# Install plugin
```
cd ..
git clone --recurse-submodules https://github.com/TauferLab/Src_ContainerizedEnv
cd Src_ContainerizedEnv
#temporary: `git checkout apptainer_plugin
go mod tidy
apptainer plugin compile plugin/
```
