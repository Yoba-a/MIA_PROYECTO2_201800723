package structs

import "container/list"

/**************************************************************
	Definicion de structs
***************************************************************/
type Mbr struct {
	Tamano int64
	Fecha  [16]byte
	Fit    [2]byte
	Firma  int8
	Tabla  [4]Particion
}

type Particion struct {
	Status byte
	Type   byte
	Fit    [2]byte
	Start  int64
	Size   int64
	Name   [16]byte
}

type Ebr struct {
	Status byte
	Fit    [2]byte
	Start  int64
	Size   int64
	Next   int64
	Name   [16]byte
}

type DiscoMontado struct {
	Path   [100]byte
	ID     int
	Estado int
	Lista  [100]ParticionMontada
}

type ParticionMontada struct {
	ID            [4]byte
	Nombre        [16]byte
	EstadoFormato byte
	EstadoMount   byte
}

type Superbloque struct {
	S_filesystem_type   int
	S_inodes_count      int64
	S_blocks_count      int64
	S_free_blocks_count int64
	S_free_inodes_count int64
	Mtime               [24]byte
	S_mnt_count         int
	S_magic             int64
	S_inode_size        int64
	S_block_size        int64
	S_first_ino         int64
	S_first_blo         int64
	S_bm_inode_start    int64
	S_bm_block_start    int64
	S_inode_start       int64
	S_block_start       int64
}

type Content struct {
	B_name  [12]byte
	B_inodo int64
}

type FolderBlock struct {
	B_contect [4]Content
}

type Fileblock struct {
	B_content [64]byte
}

type Pointer_block struct {
	B_pointers [16]int64
}

type Inodes struct {
	I_uid   int64
	I_gid   int64
	I_size  int64
	I_atime [16]byte
	I_ctime [16]byte
	I_mtime [16]byte
	I_block [15]int64
	I_type  byte
	I_perm  byte
}

type Avd struct {
	AVDFechaCreacion            [16]byte
	AVDNombreDirectorio         [16]byte
	AVDApArraySubdirectorios    [6]int64
	AVDApDetalleDirectorio      int64
	AVDApArbolVirtualDirectorio int64
	AVDProper                   [16]byte
}

type Dd struct {
	DDArrayFiles          [5]archivo
	DDApDetalleDirectorio int64
}

type Inodo struct {
	ICountInodo            int64
	ISizeArchivo           int64
	ICountBloquesAsignados int64
	IArrayBloques          [4]int64
	IApIndirecto           int64
	IIdProper              [16]byte
}

type Bloque struct {
	DBData [25]byte
}

type Bitacora struct {
	LogTipoOperacion int64
	LogTipo          int64
	LogNombre        [16]byte
	LogContenido     int64
	LogFecha         [16]byte
}

type archivo struct {
	FileNombre           [16]byte
	FileApInodo          int64
	FileDateCreacion     [16]byte
	FileDateModificacion [16]byte
}

var discos = list.New()

func Montar_disco(montado DiscoMontado) {
	discos.PushBack(montado)

}

func Discos_montados() *list.List {
	return discos
}
