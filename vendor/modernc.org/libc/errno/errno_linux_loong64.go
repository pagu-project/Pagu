// Code generated by 'ccgo errno/gen.c -crt-import-path "" -export-defines "" -export-enums "" -export-externs X -export-fields F -export-structs "" -export-typedefs "" -header -hide _OSSwapInt16,_OSSwapInt32,_OSSwapInt64 -ignore-unsupported-alignment -o errno/errno_linux_loong64.go -pkgname errno', DO NOT EDIT.

package errno

import (
	"math"
	"reflect"
	"sync/atomic"
	"unsafe"
)

var _ = math.Pi
var _ reflect.Kind
var _ atomic.Value
var _ unsafe.Pointer

const (
	E2BIG                        = 7             // errno-base.h:11:1:
	EACCES                       = 13            // errno-base.h:17:1:
	EADDRINUSE                   = 98            // errno.h:81:1:
	EADDRNOTAVAIL                = 99            // errno.h:82:1:
	EADV                         = 68            // errno.h:51:1:
	EAFNOSUPPORT                 = 97            // errno.h:80:1:
	EAGAIN                       = 11            // errno-base.h:15:1:
	EALREADY                     = 114           // errno.h:97:1:
	EBADE                        = 52            // errno.h:33:1:
	EBADF                        = 9             // errno-base.h:13:1:
	EBADFD                       = 77            // errno.h:60:1:
	EBADMSG                      = 74            // errno.h:57:1:
	EBADR                        = 53            // errno.h:34:1:
	EBADRQC                      = 56            // errno.h:37:1:
	EBADSLT                      = 57            // errno.h:38:1:
	EBFONT                       = 59            // errno.h:42:1:
	EBUSY                        = 16            // errno-base.h:20:1:
	ECANCELED                    = 125           // errno.h:109:1:
	ECHILD                       = 10            // errno-base.h:14:1:
	ECHRNG                       = 44            // errno.h:25:1:
	ECOMM                        = 70            // errno.h:53:1:
	ECONNABORTED                 = 103           // errno.h:86:1:
	ECONNREFUSED                 = 111           // errno.h:94:1:
	ECONNRESET                   = 104           // errno.h:87:1:
	EDEADLK                      = 35            // errno.h:7:1:
	EDEADLOCK                    = 35            // errno.h:40:1:
	EDESTADDRREQ                 = 89            // errno.h:72:1:
	EDOM                         = 33            // errno-base.h:37:1:
	EDOTDOT                      = 73            // errno.h:56:1:
	EDQUOT                       = 122           // errno.h:105:1:
	EEXIST                       = 17            // errno-base.h:21:1:
	EFAULT                       = 14            // errno-base.h:18:1:
	EFBIG                        = 27            // errno-base.h:31:1:
	EHOSTDOWN                    = 112           // errno.h:95:1:
	EHOSTUNREACH                 = 113           // errno.h:96:1:
	EHWPOISON                    = 133           // errno.h:121:1:
	EIDRM                        = 43            // errno.h:24:1:
	EILSEQ                       = 84            // errno.h:67:1:
	EINPROGRESS                  = 115           // errno.h:98:1:
	EINTR                        = 4             // errno-base.h:8:1:
	EINVAL                       = 22            // errno-base.h:26:1:
	EIO                          = 5             // errno-base.h:9:1:
	EISCONN                      = 106           // errno.h:89:1:
	EISDIR                       = 21            // errno-base.h:25:1:
	EISNAM                       = 120           // errno.h:103:1:
	EKEYEXPIRED                  = 127           // errno.h:111:1:
	EKEYREJECTED                 = 129           // errno.h:113:1:
	EKEYREVOKED                  = 128           // errno.h:112:1:
	EL2HLT                       = 51            // errno.h:32:1:
	EL2NSYNC                     = 45            // errno.h:26:1:
	EL3HLT                       = 46            // errno.h:27:1:
	EL3RST                       = 47            // errno.h:28:1:
	ELIBACC                      = 79            // errno.h:62:1:
	ELIBBAD                      = 80            // errno.h:63:1:
	ELIBEXEC                     = 83            // errno.h:66:1:
	ELIBMAX                      = 82            // errno.h:65:1:
	ELIBSCN                      = 81            // errno.h:64:1:
	ELNRNG                       = 48            // errno.h:29:1:
	ELOOP                        = 40            // errno.h:21:1:
	EMEDIUMTYPE                  = 124           // errno.h:108:1:
	EMFILE                       = 24            // errno-base.h:28:1:
	EMLINK                       = 31            // errno-base.h:35:1:
	EMSGSIZE                     = 90            // errno.h:73:1:
	EMULTIHOP                    = 72            // errno.h:55:1:
	ENAMETOOLONG                 = 36            // errno.h:8:1:
	ENAVAIL                      = 119           // errno.h:102:1:
	ENETDOWN                     = 100           // errno.h:83:1:
	ENETRESET                    = 102           // errno.h:85:1:
	ENETUNREACH                  = 101           // errno.h:84:1:
	ENFILE                       = 23            // errno-base.h:27:1:
	ENOANO                       = 55            // errno.h:36:1:
	ENOBUFS                      = 105           // errno.h:88:1:
	ENOCSI                       = 50            // errno.h:31:1:
	ENODATA                      = 61            // errno.h:44:1:
	ENODEV                       = 19            // errno-base.h:23:1:
	ENOENT                       = 2             // errno-base.h:6:1:
	ENOEXEC                      = 8             // errno-base.h:12:1:
	ENOKEY                       = 126           // errno.h:110:1:
	ENOLCK                       = 37            // errno.h:9:1:
	ENOLINK                      = 67            // errno.h:50:1:
	ENOMEDIUM                    = 123           // errno.h:107:1:
	ENOMEM                       = 12            // errno-base.h:16:1:
	ENOMSG                       = 42            // errno.h:23:1:
	ENONET                       = 64            // errno.h:47:1:
	ENOPKG                       = 65            // errno.h:48:1:
	ENOPROTOOPT                  = 92            // errno.h:75:1:
	ENOSPC                       = 28            // errno-base.h:32:1:
	ENOSR                        = 63            // errno.h:46:1:
	ENOSTR                       = 60            // errno.h:43:1:
	ENOSYS                       = 38            // errno.h:18:1:
	ENOTBLK                      = 15            // errno-base.h:19:1:
	ENOTCONN                     = 107           // errno.h:90:1:
	ENOTDIR                      = 20            // errno-base.h:24:1:
	ENOTEMPTY                    = 39            // errno.h:20:1:
	ENOTNAM                      = 118           // errno.h:101:1:
	ENOTRECOVERABLE              = 131           // errno.h:117:1:
	ENOTSOCK                     = 88            // errno.h:71:1:
	ENOTSUP                      = 95            // errno.h:30:1:
	ENOTTY                       = 25            // errno-base.h:29:1:
	ENOTUNIQ                     = 76            // errno.h:59:1:
	ENXIO                        = 6             // errno-base.h:10:1:
	EOPNOTSUPP                   = 95            // errno.h:78:1:
	EOVERFLOW                    = 75            // errno.h:58:1:
	EOWNERDEAD                   = 130           // errno.h:116:1:
	EPERM                        = 1             // errno-base.h:5:1:
	EPFNOSUPPORT                 = 96            // errno.h:79:1:
	EPIPE                        = 32            // errno-base.h:36:1:
	EPROTO                       = 71            // errno.h:54:1:
	EPROTONOSUPPORT              = 93            // errno.h:76:1:
	EPROTOTYPE                   = 91            // errno.h:74:1:
	ERANGE                       = 34            // errno-base.h:38:1:
	EREMCHG                      = 78            // errno.h:61:1:
	EREMOTE                      = 66            // errno.h:49:1:
	EREMOTEIO                    = 121           // errno.h:104:1:
	ERESTART                     = 85            // errno.h:68:1:
	ERFKILL                      = 132           // errno.h:119:1:
	EROFS                        = 30            // errno-base.h:34:1:
	ESHUTDOWN                    = 108           // errno.h:91:1:
	ESOCKTNOSUPPORT              = 94            // errno.h:77:1:
	ESPIPE                       = 29            // errno-base.h:33:1:
	ESRCH                        = 3             // errno-base.h:7:1:
	ESRMNT                       = 69            // errno.h:52:1:
	ESTALE                       = 116           // errno.h:99:1:
	ESTRPIPE                     = 86            // errno.h:69:1:
	ETIME                        = 62            // errno.h:45:1:
	ETIMEDOUT                    = 110           // errno.h:93:1:
	ETOOMANYREFS                 = 109           // errno.h:92:1:
	ETXTBSY                      = 26            // errno-base.h:30:1:
	EUCLEAN                      = 117           // errno.h:100:1:
	EUNATCH                      = 49            // errno.h:30:1:
	EUSERS                       = 87            // errno.h:70:1:
	EWOULDBLOCK                  = 11            // errno.h:22:1:
	EXDEV                        = 18            // errno-base.h:22:1:
	EXFULL                       = 54            // errno.h:35:1:
	X_ABILP64                    = 3             // <predefined>:377:1:
	X_ASM_GENERIC_ERRNO_BASE_H   = 0             // errno-base.h:3:1:
	X_ASM_GENERIC_ERRNO_H        = 0             // errno.h:3:1:
	X_ATFILE_SOURCE              = 1             // features.h:353:1:
	X_BITS_ERRNO_H               = 1             // errno.h:20:1:
	X_DEFAULT_SOURCE             = 1             // features.h:238:1:
	X_ERRNO_H                    = 1             // errno.h:23:1:
	X_FEATURES_H                 = 1             // features.h:19:1:
	X_FILE_OFFSET_BITS           = 64            // <builtin>:25:1:
	X_LOONGARCH_ARCH             = "loongarch64" // <predefined>:214:1:
	X_LOONGARCH_ARCH_LOONGARCH64 = 1             // <predefined>:340:1:
	X_LOONGARCH_FPSET            = 32            // <predefined>:265:1:
	X_LOONGARCH_SIM              = 3             // <predefined>:233:1:
	X_LOONGARCH_SPFPSET          = 32            // <predefined>:88:1:
	X_LOONGARCH_SZINT            = 32            // <predefined>:230:1:
	X_LOONGARCH_SZLONG           = 64            // <predefined>:388:1:
	X_LOONGARCH_SZPTR            = 64            // <predefined>:200:1:
	X_LOONGARCH_TUNE             = "la464"       // <predefined>:245:1:
	X_LOONGARCH_TUNE_LA464       = 1             // <predefined>:63:1:
	X_LP64                       = 1             // <predefined>:372:1:
	X_POSIX_C_SOURCE             = 200809        // features.h:292:1:
	X_POSIX_SOURCE               = 1             // features.h:290:1:
	X_STDC_PREDEF_H              = 1             // <predefined>:223:1:
	X_SYS_CDEFS_H                = 1             // cdefs.h:20:1:
	Linux                        = 1             // <predefined>:308:1:
	Unix                         = 1             // <predefined>:247:1:
)

type Ptrdiff_t = int64 /* <builtin>:3:26 */

type Size_t = uint64 /* <builtin>:9:23 */

type Wchar_t = int32 /* <builtin>:15:24 */

type X__int128_t = struct {
	Flo int64
	Fhi int64
} /* <builtin>:21:43 */ // must match modernc.org/mathutil.Int128
type X__uint128_t = struct {
	Flo uint64
	Fhi uint64
} /* <builtin>:22:44 */ // must match modernc.org/mathutil.Int128

type X__builtin_va_list = uintptr /* <builtin>:46:14 */
type X__float128 = float64        /* <builtin>:47:21 */

var _ int8 /* gen.c:2:13: */
