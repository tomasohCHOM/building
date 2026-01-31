#ifndef VECTOR_H
#define VECTOR_H

typedef struct vector vector;

unsigned int size(vector *vec);
unsigned int capacity(vector *vec);
unsigned int empty(vector *vec);

void init(vector *vec);
void *get(vector *vec, unsigned int index);
void *set(vector *vec, void *data, int index);
void append(vector *vec);
void pop(vector *vec);
void free(vector *vec);

#endif // VECTOR_H
