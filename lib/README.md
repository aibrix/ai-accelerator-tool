# Build from source
clang supporting at least C++17 is needed to build this project.
```bash
cd build
cmake ..
make
```
The output lib will be placed in build/lib.


# Usage
Preload this lib (through LD_PRELOAD_LIBRARY or /etc/ld.so.preload) and it will 
 hijack coresponding device management api (nvml for example) according to configurations read from gpu_mock_conf.toml file. This config file path can be configured with envvar GPU_MOCK_CONF_PATH and default to /opt/gpu_mock. 

```toml
# this ia an example config file of version 0.2.0 to explain available gpu mock configurations
version = "0.2.0"

[gpus]
## following are node-level configs
card_count = 2         # gpu card count
nvml_init_error = 9    # change nvmlInit return value, setting this will skip real init, see nvmlReturn_t below


## following are card-level configs
# config gpu card at index 0, note that the config will not be sanitized and whether the mocked scenario makes sense is up to you
[gpus.0]
remapping_pending = true
dram_ue = 7

# config gpu card at index 1, all available configs are listed below 
[gpus.1]
device_name = "A800"                                    # device name
arch = 7                                                # architecture, see nvmlDeviceArchitecture_t below
pci = "00:00"                                           # pci bus_id:device_id
uuid = "xxxxx"                                          # card uuid
link_gen = 4                                            # pcie link generation
link_width_current = 16                                 # pcie link width current
link_width_max = 16                                     # pcie link width max
nvlink_active = [true, true, true, true, true, true]    # if nvlink lane is active
remapping_failure = false                               # if there is a row-remapping failure, available for ampere and newer architecture
remapping_pending = false                               # if there is pending row-remapping, available for ampere and newer architecture 
sram_ue = 0                                             # count of uncorrectable ecc errors happened in SRAM
dram_ue = 0                                             # count of uncorrectable ecc errors happened in DRAM
dram_ce = 0                                             # count of correctable ecc errors happended in DRAM
retired_page_sbe = 0                                    # count of retired pages caused by single bit error, available for architectures before ampere
retired_page_dbe = 0                                    # count of retired pages caused by double bit error, available for architectures before ampere
retired_page_pending = false                            # if there is pending retired page
uncorrectable_agg_l1 = 0                                # aggregated count of uncorrectable ecc errors happened in L1 cache
uncorrectable_agg_l2 = 0                                # aggregated count of uncorrectable ecc errors happened in L2 cache
uncorrectable_agg_reg = 0                               # aggregated count of uncorrectable ecc errors happened in register file
crictl_xid = [62]                                       # always fire an event of type nvmlEventTypeXidCriticalError and data the configured xid

## following are ascend npu configs
[ascend-npus]
dcmi_init_error = 0                                     # change dcmiInit return value, setting this will skip real init
card_count = 1                                          # npu device count
[ascend-npus.0.0]                                       # config npu device at card 0 device 0
fault_codes = [0x80FA4E00]                              # list of fault codes to inject

```
enum references for config
```c++
/**
 * Return values for NVML API calls.
 */
typedef enum nvmlReturn_enum
{
    // cppcheck-suppress *
    NVML_SUCCESS = 0,                          //!< The operation was successful
    NVML_ERROR_UNINITIALIZED = 1,              //!< NVML was not first initialized with nvmlInit()
    NVML_ERROR_INVALID_ARGUMENT = 2,           //!< A supplied argument is invalid
    NVML_ERROR_NOT_SUPPORTED = 3,              //!< The requested operation is not available on target device
    NVML_ERROR_NO_PERMISSION = 4,              //!< The current user does not have permission for operation
    NVML_ERROR_ALREADY_INITIALIZED = 5,        //!< Deprecated: Multiple initializations are now allowed through ref counting
    NVML_ERROR_NOT_FOUND = 6,                  //!< A query to find an object was unsuccessful
    NVML_ERROR_INSUFFICIENT_SIZE = 7,          //!< An input argument is not large enough
    NVML_ERROR_INSUFFICIENT_POWER = 8,         //!< A device's external power cables are not properly attached
    NVML_ERROR_DRIVER_NOT_LOADED = 9,          //!< NVIDIA driver is not loaded
    NVML_ERROR_TIMEOUT = 10,                   //!< User provided timeout passed
    NVML_ERROR_IRQ_ISSUE = 11,                 //!< NVIDIA Kernel detected an interrupt issue with a GPU
    NVML_ERROR_LIBRARY_NOT_FOUND = 12,         //!< NVML Shared Library couldn't be found or loaded
    NVML_ERROR_FUNCTION_NOT_FOUND = 13,        //!< Local version of NVML doesn't implement this function
    NVML_ERROR_CORRUPTED_INFOROM = 14,         //!< infoROM is corrupted
    NVML_ERROR_GPU_IS_LOST = 15,               //!< The GPU has fallen off the bus or has otherwise become inaccessible
    NVML_ERROR_RESET_REQUIRED = 16,            //!< The GPU requires a reset before it can be used again
    NVML_ERROR_OPERATING_SYSTEM = 17,          //!< The GPU control device has been blocked by the operating system/cgroups
    NVML_ERROR_LIB_RM_VERSION_MISMATCH = 18,   //!< RM detects a driver/library version mismatch
    NVML_ERROR_IN_USE = 19,                    //!< An operation cannot be performed because the GPU is currently in use
    NVML_ERROR_MEMORY = 20,                    //!< Insufficient memory
    NVML_ERROR_NO_DATA = 21,                   //!< No data
    NVML_ERROR_VGPU_ECC_NOT_SUPPORTED = 22,    //!< The requested vgpu operation is not available on target device, becasue ECC is enabled
    NVML_ERROR_INSUFFICIENT_RESOURCES = 23,    //!< Ran out of critical resources, other than memory
    NVML_ERROR_FREQ_NOT_SUPPORTED = 24,        //!< Ran out of critical resources, other than memory
    NVML_ERROR_ARGUMENT_VERSION_MISMATCH = 25, //!< The provided version is invalid/unsupported
    NVML_ERROR_DEPRECATED  = 26,               //!< The requested functionality has been deprecated
    NVML_ERROR_NOT_READY = 27,                 //!< The system is not ready for the request
    NVML_ERROR_UNKNOWN = 999                   //!< An internal driver error occurred
} nvmlReturn_t;
```
```c++
/**
 * Simplified chip architecture
 */
#define NVML_DEVICE_ARCH_KEPLER    2 // Devices based on the NVIDIA Kepler architecture
#define NVML_DEVICE_ARCH_MAXWELL   3 // Devices based on the NVIDIA Maxwell architecture
#define NVML_DEVICE_ARCH_PASCAL    4 // Devices based on the NVIDIA Pascal architecture
#define NVML_DEVICE_ARCH_VOLTA     5 // Devices based on the NVIDIA Volta architecture
#define NVML_DEVICE_ARCH_TURING    6 // Devices based on the NVIDIA Turing architecture

#define NVML_DEVICE_ARCH_AMPERE    7 // Devices based on the NVIDIA Ampere architecture

#define NVML_DEVICE_ARCH_ADA       8 // Devices based on the NVIDIA Ada architecture

#define NVML_DEVICE_ARCH_HOPPER    9 // Devices based on the NVIDIA Hopper architecture

#define NVML_DEVICE_ARCH_UNKNOWN   0xffffffff // Anything else, presumably something newer
```

```c++
#define DCMI_OK 0
#define DCMI_ERROR_CODE_BASE (-8000)
#define DCMI_ERR_CODE_INVALID_PARAMETER             (DCMI_ERROR_CODE_BASE - 1)
#define DCMI_ERR_CODE_OPER_NOT_PERMITTED            (DCMI_ERROR_CODE_BASE - 2)
#define DCMI_ERR_CODE_MEM_OPERATE_FAIL              (DCMI_ERROR_CODE_BASE - 3)
#define DCMI_ERR_CODE_SECURE_FUN_FAIL               (DCMI_ERROR_CODE_BASE - 4)
#define DCMI_ERR_CODE_INNER_ERR                     (DCMI_ERROR_CODE_BASE - 5)
#define DCMI_ERR_CODE_TIME_OUT                      (DCMI_ERROR_CODE_BASE - 6)
#define DCMI_ERR_CODE_INVALID_DEVICE_ID             (DCMI_ERROR_CODE_BASE - 7)
#define DCMI_ERR_CODE_DEVICE_NOT_EXIST              (DCMI_ERROR_CODE_BASE - 8)
#define DCMI_ERR_CODE_IOCTL_FAIL                    (DCMI_ERROR_CODE_BASE - 9)
#define DCMI_ERR_CODE_SEND_MSG_FAIL                 (DCMI_ERROR_CODE_BASE - 10)
#define DCMI_ERR_CODE_RECV_MSG_FAIL                 (DCMI_ERROR_CODE_BASE - 11)
#define DCMI_ERR_CODE_NOT_REDAY                     (DCMI_ERROR_CODE_BASE - 12)
#define DCMI_ERR_CODE_NOT_SUPPORT_IN_CONTAINER      (DCMI_ERROR_CODE_BASE - 13)
#define DCMI_ERR_CODE_FILE_OPERATE_FAIL             (DCMI_ERROR_CODE_BASE - 14)
#define DCMI_ERR_CODE_RESET_FAIL                    (DCMI_ERROR_CODE_BASE - 15)
#define DCMI_ERR_CODE_ABORT_OPERATE                 (DCMI_ERROR_CODE_BASE - 16)
#define DCMI_ERR_CODE_IS_UPGRADING                  (DCMI_ERROR_CODE_BASE - 17)
#define DCMI_ERR_CODE_RESOURCE_OCCUPIED             (DCMI_ERROR_CODE_BASE - 20)
#define DCMI_ERR_CODE_PARTITION_NOT_RIGHT           (DCMI_ERROR_CODE_BASE - 22)
#define DCMI_ERR_CODE_CONFIG_INFO_NOT_EXIST         (DCMI_ERROR_CODE_BASE - 23)
#define DCMI_ERR_CODE_NOT_SUPPORT                   (DCMI_ERROR_CODE_BASE - 255)
```