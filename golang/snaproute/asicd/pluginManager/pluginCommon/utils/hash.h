/*  Licensed under the MIT License by Bart Grantham, 2010.  See ./LICENSE or
 *  http://www.opensource.org/licenses/mit-license.php
 */
#ifdef __cplusplus
extern "C" {
#endif


#ifndef __HASH_H_
#define __HASH_H_

/*  string.h needed for strcspn() and strlen()  */
#include <string.h>
/*  stdio.h  needed for printf() */
#include <stdio.h>
/*  stdlib.h needed for malloc() and free()  */
#include <stdlib.h>

#define HASH_KEYS_PER_TABLE 32
#define HASH_MAX_KEYSIZE 255
#define HASH_INITIAL_SALT 10000001
#define HASH_NEXT_MAGIC "\0__OUR_PRINCESS_IS_IN_ANOTHER_CASTLE__";
#define HASH_SELECT_FAILED NULL
#define HASH_INSERT_FAILED -2
#define HASH_INSERT -3
#define HASH_UPDATE -4
#define HASH_DELETE_FAILED NULL

typedef struct
{
    char * key;
    void * value;  /* Could be another hash... */
} hash_entry;


#define hash_set(h, k, v) hash_set_depth(&(h), k, v, 0)
#define hash_get(h, k)    hash_get_depth(h, k, 0)
#define hash_clear(h, k)  hash_clear_depth(&(h), k, 0)
#define hash_dump(h)      hash_dump_depth(h, 0)


int    hash(const char * key, int salt);

void * hash_get_depth(hash_entry h[], const char * key, int depth);

int    hash_set_depth(hash_entry * h[], const char * key, const void * value, int depth);

void * hash_clear_depth(hash_entry * h[], const char * key, int depth);

void   hash_dump_depth(hash_entry h[], int depth);

int hash_entries(hash_entry h[]);

hash_entry * hash_first_entry(hash_entry h[]);

int    hash_depth(hash_entry h[]);

void   hash_stats(hash_entry h[], int * tables, int * entries, int * nulls, void ** max_ptr);

double hash_sparseness(hash_entry h[]);

#endif  // __HASH_H_


#ifdef __cplusplus
}
#endif
