#include "common.hpp"
#include "injection.hpp"

#ifdef INJECT_DCMI

int injector::dcmi::OriginLibLoader::load_fn(const char *fnName) {
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

int dcmi_init(void) {
  auto injected = injector::conf->table["ascend_npus"]["dcmi_init_error"].value<int>();
  if (injected) {
    return injected.value();
  }

  if (injector::dcmi::ori->fns.find(__func__) == injector::dcmi::ori->fns.end()) {
    if (injector::dcmi::ori->load_fn(__func__) != 0) {
      return DCMI_ERR_CODE_INNER_ERR;
    }
  }
  return reinterpret_cast<dcmi_init_t>(
      reinterpret_cast<long long>(injector::dcmi::ori->fns.at(__func__)))();
}

int dcmi_get_card_list(int *card_num, int *card_list, int list_len) {
  auto injected = injector::conf->table["ascend_npus"]["card_count"].value<int>();
  if (injected) {
    *card_num = injected.value();
    for (auto idx = 0; idx < *card_num && idx < list_len; idx++) {
      card_list[idx] = idx;
    }
    return DCMI_OK;
  }

  if (injector::dcmi::ori->fns.find(__func__) == injector::dcmi::ori->fns.end()) {
    if (injector::dcmi::ori->load_fn(__func__) != 0) {
      return DCMI_ERR_CODE_INNER_ERR;
    }
  }
  return reinterpret_cast<dcmi_get_card_list_t>(
      reinterpret_cast<long long>(injector::dcmi::ori->fns.at(__func__)))(card_num, card_list, list_len);
}

int dcmi_get_device_id_in_card(int card_id, int *device_id_max, int *mcu_id, int *cpu_id) {
  auto injected = injector::conf->table["ascend_npus"]["card_count"].value<int>();
  if (injected) {
    *device_id_max = 1;
    return DCMI_OK;
  }

  if (injector::dcmi::ori->fns.find(__func__) == injector::dcmi::ori->fns.end()) {
    if (injector::dcmi::ori->load_fn(__func__) != 0) {
      return DCMI_ERR_CODE_INNER_ERR;
    }
  }
  return reinterpret_cast<dcmi_get_device_id_in_card_t>(
      reinterpret_cast<long long>(injector::dcmi::ori->fns.at(__func__)))(card_id, device_id_max, mcu_id, cpu_id);
}


int dcmi_get_device_errorcode_v2(int card_id, int device_id, int *error_count, unsigned int *error_code_list, unsigned int list_len) {
  auto fault_codes = injector::conf->table["ascend_npus"][card_id][device_id]["fault_codes"].as_array();
  if (fault_codes) {
    for(auto fault_code = fault_codes->begin(); fault_code != fault_codes->end(); fault_code++) {
      auto fault_code_value = fault_code->value<int>();
      if (fault_code_value) {
        if (*error_count < list_len) {
          error_code_list[*error_count] = fault_code_value.value();
          (*error_count)++;
        }
      }
    }
    return DCMI_OK;
  }

  if (injector::dcmi::ori->fns.find(__func__) == injector::dcmi::ori->fns.end()) {
    if (injector::dcmi::ori->load_fn(__func__) != 0) {
      return DCMI_ERR_CODE_INNER_ERR;
    }
  }
  return reinterpret_cast<dcmi_get_device_errorcode_v2_t>(
      reinterpret_cast<long long>(injector::dcmi::ori->fns.at(__func__)))(card_id, device_id, error_count, error_code_list, list_len);
}

#endif // INJECT_DCMI