#ifdef INJECT_DL
#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>

namespace injector {

typedef void *(*dlsym_t)(void *, const char *);

void *dlsym(void *handle, const char *symbol) {
  auto odlsym = (dlsym_t)dlvsym(RTLD_NEXT, "dlsym", "GLIBC_2.2.5");
  if (!odlsym) {
    return NULL;
  }
  odlsym = (dlsym_t)odlsym(RTLD_NEXT, "dlsym");
  if (!odlsym) {
    return NULL;
  }

  if (handle == RTLD_NEXT) {
    return odlsym(RTLD_NEXT, symbol);
  }

  void *injf = odlsym(NULL, symbol);
  if (injf) {
    return injf;
  }

  void *result = odlsym(handle, symbol);
  return result;
}

} // namespace injector

#endif // INJECT_DL