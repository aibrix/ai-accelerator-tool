#include "common.hpp"
#include "injection.hpp"

#ifdef INJECT_NVML

int injector::nvml::OriginLibLoader::load_fn(const char *fnName) {
  dlerror();
  auto fn = dlsym(RTLD_NEXT, fnName);
  if (fn != nullptr) {
    fns[std::string(fnName)] = fn;
    return 0;
  }
  logFile logf = logFile();
  if (logf.F) {
    fprintf(logf.F, "%s\n", dlerror());
  }
  return 1;
}


nvmlReturn_t nvmlInit_v2(void) {
  auto injected = injector::conf->table["gpus"]["nvml_init_error"].value<int>();
  if (injected) {
    return (nvmlReturn_t)injected.value();
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_DRIVER_NOT_LOADED;
    }
  }
  return reinterpret_cast<nvmlInit_v2_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))();
}

nvmlReturn_t nvmlInitWithFlags(unsigned int flags) {
  auto injected = injector::conf->table["gpus"]["nvml_init_error"].value<int>();
  if (injected) {
    return (nvmlReturn_t)injected.value();
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_DRIVER_NOT_LOADED;
    }
  }
  return reinterpret_cast<nvmlInitWithFlags_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(flags);
}

// nvmlReturn_t nvmlDeviceGetHandleByIndex_v2(unsigned int index, nvmlDevice_t*
// device) {
//     return NVML_SUCCESS;
// }

nvmlReturn_t nvmlDeviceGetCount_v2(unsigned int *deviceCount) {
  auto injected = injector::conf->table["gpus"]["card_count"].value<int>();
  if (injected) {
    *deviceCount = injected.value();
    return NVML_SUCCESS;
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_DRIVER_NOT_LOADED;
    }
  }
  return reinterpret_cast<nvmlDeviceGetCount_v2_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(deviceCount);
}

nvmlReturn_t nvmlDeviceGetName(nvmlDevice_t device, char *name,
                               unsigned int length) {
  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  auto injected = injector::conf->table["gpus"][std::to_string(idx)]["device_name"]
                      .value<std::string>();
  if (injected) {
    injected.value().copy(name, length);
    return NVML_SUCCESS;
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  return reinterpret_cast<nvmlDeviceGetName_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(device, name, length);
}

nvmlReturn_t nvmlDeviceGetMaxPcieLinkGeneration(nvmlDevice_t device,
                                                unsigned int *maxLinkGen) {
  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  auto injected =
      injector::conf->table["gpus"][std::to_string(idx)]["link_gen"].value<int>();
  if (injected) {
    *maxLinkGen = injected.value();
    return NVML_SUCCESS;
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  return reinterpret_cast<nvmlDeviceGetMaxPcieLinkGeneration_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(device, maxLinkGen);
}

nvmlReturn_t nvmlDeviceGetMaxPcieLinkWidth(nvmlDevice_t device,
                                           unsigned int *maxLinkWidth) {
  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  auto injected =
      injector::conf->table["gpus"][std::to_string(idx)]["link_width_max"].value<int>();
  if (injected) {
    *maxLinkWidth = injected.value();
    return NVML_SUCCESS;
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  return reinterpret_cast<nvmlDeviceGetMaxPcieLinkGeneration_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(device, maxLinkWidth);
}

nvmlReturn_t nvmlDeviceGetCurrPcieLinkWidth(nvmlDevice_t device,
                                            unsigned int *currLinkWidth) {
  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  auto injected = injector::conf->table["gpus"][std::to_string(idx)]["link_width_current"]
                      .value<int>();
  if (injected) {
    *currLinkWidth = injected.value();
    return NVML_SUCCESS;
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  return reinterpret_cast<nvmlDeviceGetMaxPcieLinkGeneration_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(device,
                                                          currLinkWidth);
}

nvmlReturn_t nvmlDeviceGetMemoryErrorCounter(nvmlDevice_t device,
                                             nvmlMemoryErrorType_t errorType,
                                             nvmlEccCounterType_t counterType,
                                             nvmlMemoryLocation_t locationType,
                                             unsigned long long *count) {
  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  if (errorType == NVML_MEMORY_ERROR_TYPE_UNCORRECTED &&
      counterType == NVML_VOLATILE_ECC &&
      locationType == NVML_MEMORY_LOCATION_SRAM) {
    auto injected = injector::conf->table["gpus"][std::to_string(idx)]["sram_ue"]
                        .value<unsigned long long>();
    if (injected) {
      *count = injected.value();
      return NVML_SUCCESS;
    }
  }
  if (errorType == NVML_MEMORY_ERROR_TYPE_UNCORRECTED &&
      counterType == NVML_VOLATILE_ECC &&
      locationType == NVML_MEMORY_LOCATION_DRAM) {
    auto injected = injector::conf->table["gpus"][std::to_string(idx)]["dram_ue"]
                        .value<unsigned long long>();
    if (injected) {
      *count = injected.value();
      return NVML_SUCCESS;
    }
  }
  if (errorType == NVML_MEMORY_ERROR_TYPE_CORRECTED &&
      counterType == NVML_VOLATILE_ECC &&
      locationType == NVML_MEMORY_LOCATION_DRAM) {
    auto injected = injector::conf->table["gpus"][std::to_string(idx)]["dram_ce"]
                        .value<unsigned long long>();
    if (injected) {
      *count = injected.value();
      return NVML_SUCCESS;
    }
  }
  if (errorType == NVML_MEMORY_ERROR_TYPE_UNCORRECTED &&
      counterType == NVML_AGGREGATE_ECC &&
      locationType == NVML_MEMORY_LOCATION_L1_CACHE) {
    auto injected =
        injector::conf->table["gpus"][std::to_string(idx)]["uncorrectable_agg_l1"]
            .value<unsigned long long>();
    if (injected) {
      *count = injected.value();
      return NVML_SUCCESS;
    }
  }
  if (errorType == NVML_MEMORY_ERROR_TYPE_UNCORRECTED &&
      counterType == NVML_AGGREGATE_ECC &&
      locationType == NVML_MEMORY_LOCATION_L2_CACHE) {
    auto injected =
        injector::conf->table["gpus"][std::to_string(idx)]["uncorrectable_agg_l2"]
            .value<unsigned long long>();
    if (injected) {
      *count = injected.value();
      return NVML_SUCCESS;
    }
  }
  if (errorType == NVML_MEMORY_ERROR_TYPE_UNCORRECTED &&
      counterType == NVML_AGGREGATE_ECC &&
      locationType == NVML_MEMORY_LOCATION_REGISTER_FILE) {
    auto injected =
        injector::conf->table["gpus"][std::to_string(idx)]["uncorrectable_agg_reg"]
            .value<unsigned long long>();
    if (injected) {
      *count = injected.value();
      return NVML_SUCCESS;
    }
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  return reinterpret_cast<nvmlDeviceGetMemoryErrorCounter_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(
      device, errorType, counterType, locationType, count);
}

nvmlReturn_t nvmlDeviceGetRemappedRows(nvmlDevice_t device,
                                       unsigned int *corrRows,
                                       unsigned int *uncRows,
                                       unsigned int *isPending,
                                       unsigned int *failureOccurred) {
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  auto oret = reinterpret_cast<nvmlDeviceGetRemappedRows_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(
      device, corrRows, uncRows, isPending, failureOccurred);

  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  auto injected = injector::conf->table["gpus"][std::to_string(idx)]["remapping_failure"]
                      .value<bool>();
  bool flagModified = false;
  if (injected) {
    if (injected.value()) {
      *failureOccurred = 1;
    } else {
      *failureOccurred = 0;
    }
    flagModified = true;
  }
  injected = injector::conf->table["gpus"][std::to_string(idx)]["remapping_pending"]
                 .value<bool>();
  if (injected) {
    if (injected.value()) {
      *isPending = 1;
    } else {
      *isPending = 0;
    }
    flagModified = true;
  }
  if (flagModified)
    return NVML_SUCCESS;
  return oret;
}

nvmlReturn_t nvmlDeviceGetRetiredPages(nvmlDevice_t device,
                                       nvmlPageRetirementCause_t cause,
                                       unsigned int *pageCount,
                                       unsigned long long *addresses) {
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  auto oret = reinterpret_cast<nvmlDeviceGetRetiredPages_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(
      device, cause, pageCount,
      addresses); // should not call origin after changing pageCount or it may
                  // overflow addresses

  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }

  if (cause == NVML_PAGE_RETIREMENT_CAUSE_MULTIPLE_SINGLE_BIT_ECC_ERRORS) {
    auto injected = injector::conf->table["gpus"][std::to_string(idx)]["retired_page_sbe"]
                        .value<int>();
    if (injected) {
      *pageCount = injected.value();
      return NVML_SUCCESS;
    }
  }
  if (cause == NVML_PAGE_RETIREMENT_CAUSE_DOUBLE_BIT_ECC_ERROR) {
    auto injected = injector::conf->table["gpus"][std::to_string(idx)]["retired_page_dbe"]
                        .value<int>();
    if (injected) {
      *pageCount = injected.value();
      return NVML_SUCCESS;
    }
  }
  return oret;
}

nvmlReturn_t
nvmlDeviceGetRetiredPagesPendingStatus(nvmlDevice_t device,
                                       nvmlEnableState_t *isPending) {
  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  auto injected =
      injector::conf->table["gpus"][std::to_string(idx)]["retired_page_pending"]
          .value<bool>();
  if (injected) {
    if (injected.value()) {
      *isPending = NVML_FEATURE_ENABLED;
    } else {
      *isPending = NVML_FEATURE_DISABLED;
    }
    return NVML_SUCCESS;
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  return reinterpret_cast<nvmlDeviceGetRetiredPagesPendingStatus_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(device, isPending);
}

nvmlReturn_t nvmlDeviceGetArchitecture(nvmlDevice_t device,
                                       nvmlDeviceArchitecture_t *arch) {
  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  auto injected = injector::conf->table["gpus"][std::to_string(idx)]["arch"].value<int>();
  if (injected) {
    *arch = (nvmlDeviceArchitecture_t)injected.value();
    return NVML_SUCCESS;
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  return reinterpret_cast<nvmlDeviceGetArchitecture_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(device, arch);
}

nvmlReturn_t nvmlDeviceGetUUID(nvmlDevice_t device, char *uuid,
                               unsigned int length) {
  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  auto injected =
      injector::conf->table["gpus"][std::to_string(idx)]["uuid"].value<std::string>();
  if (injected) {
    injected.value().copy(uuid, length);
    return NVML_SUCCESS;
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  return reinterpret_cast<nvmlDeviceGetUUID_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(device, uuid, length);
}

nvmlReturn_t nvmlDeviceGetNvLinkState(nvmlDevice_t device, unsigned int link,
                                      nvmlEnableState_t *isActive) {
  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  auto injected =
      injector::conf->table["gpus"][std::to_string(idx)]["nvlink_active"][link]
          .value<bool>();
  if (injected) {
    if (injected.value()) {
      *isActive = NVML_FEATURE_ENABLED;
    } else {
      *isActive = NVML_FEATURE_DISABLED;
    }
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  return reinterpret_cast<nvmlDeviceGetNvLinkState_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(device, link,
                                                          isActive);
}

nvmlReturn_t nvmlDeviceGetFieldValues(nvmlDevice_t device, int valuesCount,
                                      nvmlFieldValue_t *values) {
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  auto oret = reinterpret_cast<nvmlDeviceGetFieldValues_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(device, valuesCount,
                                                          values);

  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  auto injected =
      injector::conf->table["gpus"][std::to_string(idx)]["nvlink_active"].as_array();
  if (injected) {
    for (int i = 0; i < valuesCount; ++i) {
      if (values[i].fieldId == NVML_FI_DEV_NVLINK_LINK_COUNT) {
        values[i].valueType = NVML_VALUE_TYPE_UNSIGNED_INT;
        values[i].value = nvmlValue_t{.uiVal = (unsigned int)injected->size()};
        oret = NVML_SUCCESS;
      }
    }
  }
  return oret;
}

nvmlReturn_t nvmlDeviceGetPciInfo_v3(nvmlDevice_t device, nvmlPciInfo_t *pci) {
  unsigned int idx = 0;
  if (auto ret = nvmlDeviceGetIndex(device, &idx); ret != NVML_SUCCESS) {
    return ret;
  }
  auto injected =
      injector::conf->table["gpus"][std::to_string(idx)]["pci"].value<std::string>();
  if (injected) {
    auto pciConfig = injected.value();
    pci->bus = std::stol(injector::find_nth(pciConfig, ":", 0), NULL, 16);
    pci->device = std::stol(injector::find_nth(pciConfig, ":", 1), NULL, 16);
    return NVML_SUCCESS;
  }
  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  return reinterpret_cast<nvmlDeviceGetPciInfo_v3_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(device, pci);
}

nvmlReturn_t nvmlEventSetWait_v2(nvmlEventSet_t set, nvmlEventData_t *data,
                                 unsigned int timeoutms) {
  bool inject_flag = true;

  unsigned int device_count;
  if (auto injected = injector::conf->table["gpus"]["card_count"].value<int>();
      injected) {
    device_count = injected.value();
  } else {
    if (auto ret = nvmlDeviceGetCount(&device_count); ret != NVML_SUCCESS) {
      inject_flag = false;
    }
  }

  if (inject_flag) {
    for (unsigned int idx = 0; idx < device_count; ++idx) {
      if (auto injected =
              injector::conf->table["gpus"][std::to_string(idx)]["crictl_xid"][0]
                  .value<int>();
          injected) {
        // current only use first configed xid
        auto xid = injected.value();
        // we do not know how to consume nvmlEventSet_t since we do not know its
        // definition, therefore we can not access the registered devices and
        // event types
        nvmlDevice_t device;
        if (auto ret = nvmlDeviceGetHandleByIndex(idx, &device);
            ret != NVML_SUCCESS) {
          continue;
        }
        data->device = device;
        data->eventType = nvmlEventTypeXidCriticalError;
        data->eventData = xid;
        return NVML_SUCCESS;
      }
    }
  }

  if (injector::nvml::ori->fns.find(__func__) == injector::nvml::ori->fns.end()) {
    if (injector::nvml::ori->load_fn(__func__) != 0) {
      return NVML_ERROR_FUNCTION_NOT_FOUND;
    }
  }
  return reinterpret_cast<nvmlEventSetWait_v2_t>(
      reinterpret_cast<long long>(injector::nvml::ori->fns.at(__func__)))(set, data, timeoutms);
}

#endif // INJECT_NVML