#pragma once
#include "nvml.h"


typedef nvmlReturn_t (*nvmlInit_v2_t)(void);
typedef nvmlReturn_t (*nvmlInitWithFlags_t)(unsigned int);
typedef nvmlReturn_t (*nvmlDeviceGetCount_v2_t)(unsigned int *);
typedef nvmlReturn_t (*nvmlDeviceGetName_t)(nvmlDevice_t, char *, unsigned int);
typedef nvmlReturn_t (*nvmlDeviceGetMaxPcieLinkGeneration_t)(nvmlDevice_t,
                                                             unsigned int *);
typedef nvmlReturn_t (*nvmlDeviceGetMaxPcieLinkWidth_t)(nvmlDevice_t,
                                                        unsigned int *);
typedef nvmlReturn_t (*nvmlDeviceGetCurrPcieLinkWidth_t)(nvmlDevice_t,
                                                         unsigned int *);
typedef nvmlReturn_t (*nvmlDeviceGetMemoryErrorCounter_t)(nvmlDevice_t,
                                                          nvmlMemoryErrorType_t,
                                                          nvmlEccCounterType_t,
                                                          nvmlMemoryLocation_t,
                                                          unsigned long long *);
typedef nvmlReturn_t (*nvmlDeviceGetRemappedRows_t)(nvmlDevice_t,
                                                    unsigned int *,
                                                    unsigned int *,
                                                    unsigned int *,
                                                    unsigned int *);
typedef nvmlReturn_t (*nvmlDeviceGetRetiredPages_t)(nvmlDevice_t,
                                                    nvmlPageRetirementCause_t,
                                                    unsigned int *,
                                                    unsigned long long *);
typedef nvmlReturn_t (*nvmlDeviceGetRetiredPagesPendingStatus_t)(
    nvmlDevice_t, nvmlEnableState_t *);
typedef nvmlReturn_t (*nvmlDeviceGetArchitecture_t)(nvmlDevice_t,
                                                    nvmlDeviceArchitecture_t *);
typedef nvmlReturn_t (*nvmlDeviceGetUUID_t)(nvmlDevice_t, char *, unsigned int);
typedef nvmlReturn_t (*nvmlDeviceGetNvLinkState_t)(nvmlDevice_t, unsigned int,
                                                   nvmlEnableState_t *);
typedef nvmlReturn_t (*nvmlDeviceGetFieldValues_t)(nvmlDevice_t, int,
                                                   nvmlFieldValue_t *);
typedef nvmlReturn_t (*nvmlDeviceGetPciInfo_v3_t)(nvmlDevice_t,
                                                  nvmlPciInfo_t *);
typedef nvmlReturn_t (*nvmlDeviceGetIndex_t)(nvmlDevice_t, unsigned int *);
typedef nvmlReturn_t (*nvmlEventSetWait_v2_t)(nvmlEventSet_t, nvmlEventData_t *,
                                              unsigned int);