#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "libimgutil.h"

int main(int argc, char** argv) {
  char *path, *err;
  GoUint w, h;

  if (argc < 2) {
    fprintf(stderr, "missing argument\n");
    return 1;
  }

  path = strdup(argv[1]);
  err = ImgutilGetImageSize(path, &w, &h);
  free(path);

  if (err != NULL) {
    fprintf(stderr, "error: %s\n", err);
    free(err);
    return 1;
  }

  printf("%s: %llux%llu\n", argv[1], w, h);

  return 0;
}