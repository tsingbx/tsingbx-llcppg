#include <stddef.h>
/* The error value type gpg_error_t.  */

/* We would really like to use bit-fields in a struct, but using
 * structs as return values can cause binary compatibility issues, in
 * particular if you want to do it efficiently (also see
 * -freg-struct-return option to GCC).  */
typedef unsigned int gpg_error_t;
/* String functions.  */

/* Return a pointer to a string containing a description of the error
 * code in the error value ERR.  This function is not thread-safe.  */
const char *gpg_strerror (gpg_error_t err);

/* Return the error string for ERR in the user-supplied buffer BUF of
 * size BUFLEN.  This function is, in contrast to gpg_strerror,
 * thread-safe if a thread-safe strerror_r() function is provided by
 * the system.  If the function succeeds, 0 is returned and BUF
 * contains the string describing the error.  If the buffer was not
 * large enough, ERANGE is returned and BUF contains as much of the
 * beginning of the error string as fits into the buffer.  */
int gpg_strerror_r (gpg_error_t err, char *buf, size_t buflen);

/* Return a pointer to a string containing a description of the error
 * source in the error value ERR.  */
const char *gpg_strsource (gpg_error_t err);

/* Only use free slots, never change or reorder the existing
 * entries.  */
typedef enum
{
    GPG_ERR_NO_ERROR = 0,
    GPG_ERR_GENERAL = 1,
    GPG_ERR_UNKNOWN_PACKET = 2,
    /* This is one more than the largest allowed entry.  */
    GPG_ERR_CODE_DIM = 65536
} gpg_err_code_t;

typedef struct
{
  long _vers;
  union {
    volatile char _priv[64];
    long _x_align;
    long *_xp_align;
  } u;
} gpgrt_lock_t;



/* NB: If GPGRT_LOCK_DEFINE is not used, zero out the lock variable
   before passing it to gpgrt_lock_init.  */
gpg_err_code_t gpgrt_lock_init (gpgrt_lock_t *lockhd);
gpg_err_code_t gpgrt_lock_lock (gpgrt_lock_t *lockhd);
gpg_err_code_t gpgrt_lock_trylock (gpgrt_lock_t *lockhd);
gpg_err_code_t gpgrt_lock_unlock (gpgrt_lock_t *lockhd);
gpg_err_code_t gpgrt_lock_destroy (gpgrt_lock_t *lockhd);
