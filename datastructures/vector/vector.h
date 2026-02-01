#ifndef VECTOR_H
#define VECTOR_H

#include <stdint.h>

typedef struct vector vector;

uint32_t size(vector *vec);
uint32_t capacity(vector *vec);
uint32_t empty(vector *vec);

vector *init_vector(vector *vec, uint32_t capacity);
void append(vector *vec, int elem);
int pop(vector *vec);
void free_vec(vector *vec);

int get(vector *vec, uint32_t index);
void set(vector *vec, uint32_t index, int elem);

#endif // VECTOR_H
