#include "vector.h"

struct vector {
  unsigned int size;
  unsigned int capacity;
};

unsigned int size(vector *vec) { return vec->size; }

unsigned int capacity(vector *vec) { return vec->capacity; }

unsigned int empty(vector *vec) { return size(vec) == 0; }
