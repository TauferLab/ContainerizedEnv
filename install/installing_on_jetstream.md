#   Remove go
First, check if go is installed `which go`. If so, delete the binary
```
sudo rm -rf /path/to/go
```

# Install go=1.19

You must first install development tools and libraries to your host.
On Debian-based systems, including Ubuntu:

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

First, download the Go tar.gz archive to `/tmp`, then extract the archive to `/usr/local`.
```
sudo apt install uidmap
wget -O /tmp/go1.19.linux-amd64.tar.gz https://dl.google.com/go/go1.19.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf /tmp/go1.19.linux-amd64.tar.gz
```

Finally, add /usr/local/go/bin to the PATH environment variable in your bashrc:
```
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

# Install apptainer
```
git clone https://github.com/apptainer/apptainer.git
cd apptainer
```

Use version 1.1.0, and then run the following commands to install Apptainer
```
git checkout v1.1.0

./mconfig
cd ./builddir
make
sudo make install
```
Check apptainer version: `apptainer --version`

# Install plugin

Clone the ContainerizedEnv repository inside the apptainer source code tree (repository)
```
cd /path/to/apptainer
git clone --recurse-submodules https://github.com/TauferLab/ContainerizedEnv
```
Remove go.mod and go.sum files that are stored inside the plugin folder
```
cd /path/to/ContainerizedEnv/plugin/
rm -rf go.mod go.sum
```

Manage go dependencies
```
cd /path/to/ContainerizedEnv
go mod tidy
```

Compile and install plugin
```
apptainer plugin compile plugin/
sudo apptainer plugin install plugin/plugin.sif
```