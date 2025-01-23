#ifdef INJECT_DL
#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>

namespace injector {

typedef void *(*dlsym_t)(void *, const char *);

// Rename to real_dlsym to avoid recursion
extern "C" void *real_dlsym(void *handle, const char *symbol) {
#ifdef __linux__
  auto odlsym = (dlsym_t)dlvsym(RTLD_NEXT, "dlsym", "GLIBC_2.2.5");
#else
  // On macOS, we need to get the next dlsym in the chain
  static void *(*next_dlsym)(void *, const char *) = NULL;
  if (!next_dlsym) {
    // Get handle to next loaded library
    void *handle = dlopen("/usr/lib/libSystem.B.dylib", RTLD_LAZY);
    if (handle) {
      next_dlsym = (void *(*)(void *, const char *))dlsym(handle, "dlsym");
    }
  }
  auto odlsym = next_dlsym;
#endif
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