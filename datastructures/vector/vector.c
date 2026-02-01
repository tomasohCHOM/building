#include "vector.h"
#include <stdint.h>
#include <stdlib.h>

struct vector {
  int *data;
  uint32_t size;
  uint32_t capacity;
};

unsigned int size(vector *vec) { return vec->size; }

unsigned int capacity(vector *vec) { return vec->capacity; }

unsigned int empty(vector *vec) { return size(vec) == 0; }

vector *init_vector(vector *vec, uint32_t capacity) {
  vec->data = malloc(capacity * sizeof(int));
  if (!vec->data) {
    return NULL;
  }
  vec->capacity = capacity;
  vec->size = 0;
  return vec;
}

void append(vector *vec, int elem) {
  if (vec->size == vec->capacity) {
    uint32_t new_capacity = vec->capacity * 2;
    int *new_data = realloc(vec->data, new_capacity * sizeof(int));
    if (!new_data) {
      return;
    }
    vec->data = new_data;
    vec->capacity = new_capacity;
  }
  vec->data[vec->size++] = elem;
}

int pop(vector *vec) {
  int elem = vec->data[vec->size - 1];
  vec->size--;
  return elem;
}

int get(vector *vec, uint32_t index) {
  if (index >= vec->size) {
    return -1;
  }
  return vec->data[index];
}

void set(vector *vec, uint32_t index, int elem) {
  if (index >= vec->size) {
    return;
  }
  vec->data[index] = elem;
}
