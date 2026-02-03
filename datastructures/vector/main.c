#include "vector.h"
#include <stdint.h>
#include <stdio.h>

int main() {
  vector *vec = init_vector(4);
  printf("capacity: %d\n", capacity(vec));
  printf("size before adding elements: %d\n", size(vec));

  for (uint32_t i = 0; i < 5; i++)
    append(vec, i);

  printf("size after adding elements: %d\n", size(vec));
  printf("%d\n", get(vec, 2));
  return 0;
}
