# AI Accelerator Tool

A command-line diagnostic tool for GPU health monitoring and troubleshooting. This tool helps identify and diagnose common GPU issues, including memory leaks, hardware failures, and performance degradation.

# Features

- Real-time GPU health monitoring
- Memory leak detection
- Hardware failure diagnosis
- Performance metrics analysis
- Mock testing capabilities for development

# Run Tests

Run tests in docker container:
```bash
make test
```

Run tests locally and generate coverage report:

```bash
make test-local
```

# Build

## Build gpu injection library (Only supports compilation in Linux environment)

If you are developing on MacOS, you can consider using a docker container for compilation. 
Taking the `ubuntu:22.04` image as an example, you need to install the following dependencies in the container and mount the project into the container for compilation.

1. Start the container

```bash
docker run --platform=linux/amd64 -itd -v ./ai-accelerator-tool:/git/src/github.com/aibrix/ai-accelerator-tool/ ubuntu:22.04
```

2. Install dependencies in the container

```bash
apt update
apt install -y vim cmake clang libnvidia-ml-dev git wget
wget https://go.dev/dl/go1.23.2.linux-amd64.tar.gz
tar xvf go1.23.2.linux-amd64.tar.gz
echo "export PATH=$PATH:/go/bin" >> ~/.bashrc
source ~/.bashrc
```

3. Compile the project in the container

```bash
cd /git/src/github.com/aibrix/ai-accelerator-tool
git submodule update --init --recursive
make lib-injection
cp lib/build/lib/libdevso-injection.so pkg/mock/resources/injectiond.so
```

## Build gpu-diagnosis

```bash
GOOS=linux GOARCH=amd64 make build
```

The binary will be generated in `bin/`.

# Usage

## GPU Diagnosis

```bash

# Set the number of GPU cards in the machine, for example, 4.
export GPU_CARD_COUNT=4

# Run the diagnosis.
ai-accelerator-tool diagnose
```

Note:
- This tool requires the `nvidia-smi` command to be installed.

## GPU Exception Mock

### 1. Prepare configuration files for fault simulation.

You can refer to the comments in `hack/gpu_mock_conf.toml` to configure the fault scenario.

### 2. Prepare shared library for fault simulation.

#### Method 1: Use shared library through ai-accelerator-tool.

```bash
ai-accelerator-tool mock --config /PATH/TO/gpu_mock_conf.toml
```

#### Method 2: Use shared library manually.

```bash
mkdir -p /opt/gpu_mock && cd /opt/gpu_mock/

cp /PATH/TO/nvml_injectiond.so /opt/gpu_mock/
cp /PATH/TO/gpu_mock_conf.toml /opt/gpu_mock/

echo "/opt/gpu_mock/nvml_injectiond.so" >> /etc/ld.so.preload
```
