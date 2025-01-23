#pragma once

#include <dlfcn.h>
#include <string>
#include <unordered_map>
#define TOML_EXCEPTIONS 0
#include <sys/stat.h>
#include <toml++/toml.hpp>

#define CONTAINER_HOST_MOUNT_PATH "/host"
#define DEFAULT_CONF_PATH "/opt/gpu_mock/gpu_mock_conf.toml"
#define DEFAULT_FOLDER "/opt/gpu_mock"
#define LOG_NAME "/log"
#define concat(a, b) a b

namespace injector {

const auto CONF_PATH_ENV = "GPU_MOCK_CONF_PATH";

typedef struct {
  bool success;
  std::string err_msg;
  toml::table table;
} conf_t;

inline bool exist_file(const char *path) {
  struct stat buf;
  return stat(path, &buf) == 0;
}

inline bool exist_folder(const char *path) {
  struct stat buf;
  return stat(path, &buf) == 0 && !(buf.st_mode & S_IFDIR);
}

inline std::string find_nth(std::string s, std::string sep, int n) {
  auto start = 0U;
  for (auto end = s.find(sep); end != std::string::npos;
       end = s.find(sep, start)) {
    if (n == 0) {
      return s.substr(start, end - start);
    }
    n--;
    start = end + sep.size();
  }
  return s.substr(start);
}

inline std::string search_conf(bool *success) {
  if (auto env_path = std::getenv(CONF_PATH_ENV); env_path != nullptr) {
    if (exist_file(env_path)) {
      return env_path;
    }
  }
  if (exist_file(concat(CONTAINER_HOST_MOUNT_PATH, DEFAULT_CONF_PATH))) {
    return concat(CONTAINER_HOST_MOUNT_PATH, DEFAULT_CONF_PATH);
  }
  if (exist_file(DEFAULT_CONF_PATH)) {
    return DEFAULT_CONF_PATH;
  }
  *success = false;
  return "";
}

inline conf_t *load_conf() {
  auto conf = new conf_t;
  bool has_conf = true;
  auto conf_path = search_conf(&has_conf);
  if (!has_conf) {
    conf->success = false;
    conf->err_msg = "failed to find gpu_mock_conf.toml";
    return conf;
  }
  auto result = toml::parse_file(conf_path);
  if (!result) {
    conf->success = false;
    char buf[255];
    auto parse_err_desp = result.error().description();
    auto parse_err_pos = result.error().source().begin;
    std::string err_path;
    if (result.error().source().path) {
      err_path = *result.error().source().path;
    }
    std::snprintf(
        buf, 255, "toml parse error at line %d column %d of %s: %.*s",
        parse_err_pos.line, parse_err_pos.column, err_path.c_str(),
        static_cast<int>(parse_err_desp.length()),
        parse_err_desp.data()); // todo: get rid of this string_view hack
    conf->err_msg = buf;
    return conf;
  }
  conf->success = true;
  conf->table = std::move(result).table();
  return conf;
}

class logFile {
public:
  FILE *F;
  logFile() {
    if (exist_folder(concat(CONTAINER_HOST_MOUNT_PATH, DEFAULT_FOLDER))) {
      F = fopen(
          concat(concat(CONTAINER_HOST_MOUNT_PATH, DEFAULT_FOLDER), LOG_NAME),
          "a");
    } else if (exist_folder(DEFAULT_FOLDER)) {
      F = fopen(concat(DEFAULT_FOLDER, LOG_NAME), "a");
    }
  }
  ~logFile() {
    if (F) {
      fclose(F);
    }
  }
};

static conf_t *conf;

#ifdef INJECT_NVML
namespace nvml {
class OriginLibLoader {
public:
  OriginLibLoader() =default;
  ~OriginLibLoader() = default;
  int load_fn(const char *fnName);
  std::unordered_map<std::string, void *> fns;
};

static OriginLibLoader *ori = nullptr;
} // namespace nvml
#endif // INJECT_NVML

#ifdef INJECT_DCMI
namespace dcmi {
class OriginLibLoader {
public:
  OriginLibLoader() = default;
  ~OriginLibLoader() = default;
  int load_fn(const char *fnName);
  std::unordered_map<std::string, void *> fns;
};

static OriginLibLoader *ori = nullptr;
}
#endif // INJECT_DCMI

static void __attribute__((constructor)) init(void) {
#ifdef INJECT_NVML
  if (nvml::ori == nullptr) {
    nvml::ori = new nvml::OriginLibLoader;
  }
#endif
#ifdef INJECT_DCMI
  if (dcmi::ori == nullptr) {
    dcmi::ori = new dcmi::OriginLibLoader;
  }
#endif
  if (conf == nullptr) {
    conf = load_conf();
  }
  if (conf != nullptr) {
    if (!conf->success) {
      logFile logf = logFile();
      if (logf.F) {
        fprintf(logf.F, "%s\n", conf->err_msg.c_str());
      }
    }
  }
}

static void __attribute__((destructor)) destory(void) {
#ifdef INJECT_NVML
  if (nvml::ori != nullptr) {
    delete nvml::ori;
  }
#endif
#ifdef INJECT_DCMI
  if (dcmi::ori != nullptr) {
    delete dcmi::ori;
  }
#endif
  if (conf != nullptr) {
    delete conf;
  }
}

} // namespace injector