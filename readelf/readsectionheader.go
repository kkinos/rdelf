package readelf

import (
	"bytes"
	"encoding/binary"

	"unsafe"
)

type sectionHeader struct {
	Sh_name      Elf64_Word  /* Section name (string tbl index) */
	Sh_type      Elf64_Word  /* Section type */
	Sh_flags     Elf64_Xword /* Section flags */
	Sh_addr      Elf64_Addr  /* Section virtual addr at execution */
	Sh_offset    Elf64_Off   /* Section file offset */
	Sh_size      Elf64_Xword /* Section size in bytes */
	Sh_link      Elf64_Word  /* Link to another section */
	Sh_info      Elf64_Word  /* Additional section information */
	Sh_addralign Elf64_Xword /* Section alignment */
	Sh_entsize   Elf64_Xword /* Entry size if section holds table */
}

type SectionHeaderInfo struct {
	Name      string
	Type      string
	Size      Elf64_Xword
	EntrySize Elf64_Xword
}

type SectionHeaderInfos []SectionHeaderInfo

type SHType uint32

const (
	SHT_NULL          SHType = 0
	SHT_PROGBITS      SHType = 1
	SHT_SYMTAB        SHType = 2
	SHT_STRTAB        SHType = 3
	SHT_RELA          SHType = 4
	SHT_HASH          SHType = 5
	SHT_DYNAMIC       SHType = 6
	SHT_NOTE          SHType = 7
	SHT_NOBITS        SHType = 8
	SHT_REL           SHType = 9
	SHT_SHLIB         SHType = 10
	SHT_DYNSYM        SHType = 11
	SHT_INIT_ARRAY    SHType = 14
	SHT_FINI_ARRAY    SHType = 15
	SHT_PREINIT_ARRAY SHType = 16
	SHT_GROUP         SHType = 17
	SHT_SYMTAB_SHNDX  SHType = 18
	SHT_NUM           SHType = 19
	SHT_GNU_HASH      SHType = 0x6ffffff6
	SHT_GNU_verneed   SHType = 0x6ffffffe
	SHT_GNU_versym    SHType = 0x6fffffff
)

func setSHType(info SectionHeaderInfo, t SHType) SectionHeaderInfo {
	switch t {
	case SHT_NULL:
		info.Type = "Null"
	case SHT_PROGBITS:
		info.Type = "Progbit"
	case SHT_SYMTAB:
		info.Type = "Symtab"
	case SHT_STRTAB:
		info.Type = "Strtab"
	case SHT_RELA:
		info.Type = "Rela"
	case SHT_HASH:
		info.Type = "Hash"
	case SHT_DYNAMIC:
		info.Type = "Dynamic"
	case SHT_NOTE:
		info.Type = "Note"
	case SHT_NOBITS:
		info.Type = "Nobits"
	case SHT_REL:
		info.Type = "Rel"
	case SHT_SHLIB:
		info.Type = "Shlib"
	case SHT_DYNSYM:
		info.Type = "Dynsym"
	case SHT_INIT_ARRAY:
		info.Type = "Initarray"
	case SHT_FINI_ARRAY:
		info.Type = "Finiarray"
	case SHT_PREINIT_ARRAY:
		info.Type = "Preinitarray"
	case SHT_GROUP:
		info.Type = "Group"
	case SHT_SYMTAB_SHNDX:
		info.Type = "Symtabshndx"
	case SHT_NUM:
		info.Type = "Num"
	case SHT_GNU_HASH:
		info.Type = "GNU_hash"
	case SHT_GNU_verneed:
		info.Type = "GNU_verneed"
	case SHT_GNU_versym:
		info.Type = "GNU_versym"
	default:
		info.Type = "Unknown"
	}
	return info
}

func ReadSectionHeaders(file []byte, shoff Elf64_Off, shnum Elf64_Half, shsize Elf64_Half) (SectionHeaderInfos, error) {
	var infos SectionHeaderInfos
	eheader, err := ReadELFHeader(file)
	if err != nil {
		return infos, err
	}

	var order binary.ByteOrder
	if eheader.Data == "little endian" {
		order = binary.LittleEndian
	} else {
		order = binary.BigEndian
	}

	for i := 0; i < int(shnum); i++ {
		var sheader sectionHeader
		var info SectionHeaderInfo
		shs := unsafe.Sizeof(sheader)
		sh := make([]byte, shs)
		copy(sh, file[int(shoff)+int(shsize)*i:])
		shr := bytes.NewReader(sh)
		err = binary.Read(shr, order, &sheader)
		if err != nil {
			return infos, err
		}
		info = setSHType(info, SHType(sheader.Sh_type))

		infos = append(infos, info)
	}
	return infos, nil
}
