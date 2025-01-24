#pragma once

#ifdef INJECT_DCMI

#include "dcmi_interface_api.h"

typedef int (*dcmi_init_t)(void);
typedef int (*dcmi_get_card_list_t)(int *, int *, int);
typedef int (*dcmi_get_device_id_in_card_t)(int, int *, int *, int *);
typedef int (*dcmi_get_device_errorcode_v2_t)(int, int, int *, unsigned int *,
                                              unsigned int);
typedef int (*dcmi_get_device_pcie_info_v2_t)(int, int,
                                              struct dcmi_pcie_info_all *);

#endif // INJECT_DCMI