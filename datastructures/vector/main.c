#include "vector.h"
#include <stdint.h>
#include <stdio.h>

int main() {
  vector *vec;
  init_vector(vec, 4);
  for (uint32_t i = 0; i < 5; i++)
    append(vec, 5);
  printf("%d\n", get(vec, 2));
  return 0;
}
