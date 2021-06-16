package partition

import (
	"fmt"
	"io"

	"unicode/utf16"

	"github.com/rekby/gpt"
	"github.com/threez/go-sparsebundle"
)

const SectorSize = 512

type Partition gpt.Partition

// Partitions returns all partitions of the passed bundle
func Partitions(b *sparsebundle.Bundle) ([]*Partition, error) {
	_, err := b.Seek(SectorSize, io.SeekStart)
	if err != nil {
		return nil, err
	}

	tab, err := gpt.ReadTable(b, SectorSize)
	if err != nil {
		return nil, err
	}

	var partitions []*Partition
	for _, p := range tab.Partitions {
		if p.FirstLBA == 0 && p.LastLBA == 0 {
			// ignore
			continue
		}
		np := new(Partition)
		*np = Partition(p)
		partitions = append(partitions, np)
	}

	return partitions, nil
}

func (p Partition) TypeName() string {
	t, ok := PartitionTypes[p.Type.String()]
	if !ok {
		t = "unknown"
	}
	return t
}

func (p Partition) Name() string {
	var n []rune

	for i := 0; i < 72/2; i += 2 {
		if p.PartNameUTF16[i] == 0x00 {
			break
		}
		r := utf16.Decode([]uint16{
			uint16(p.PartNameUTF16[i]) + uint16(p.PartNameUTF16[i+1])<<8,
		})
		n = append(n, r...)
	}

	return string(n)
}

func (p Partition) Size() uint64 {
	return (p.LastLBA - p.FirstLBA) * SectorSize
}

func (p Partition) String() string {
	return fmt.Sprintf("<Partition type=%q name=%q First=%d Last=%d Size=%d/>",
		p.TypeName(), p.Name(), p.FirstLBA*SectorSize, p.LastLBA*SectorSize, p.Size())
}

var PartitionTypes = map[string]string{
	"EBD0A0A2-B9E5-4433-87C0-68B6B72699C7": "FAT",
	"E3C9E316-0B5C-4DB8-817D-F92DF00215AE": "FAT32",
	"DE94BBA4-06D1-4D40-A16A-BFD50179D6AC": "Windows Recovery",
	"C91818F9-8025-47AF-89D2-F030D7000C2C": "Plan 9",
	"AF9B60A0-1431-4F62-BC68-3311714A69AD": "Windows LDM data",
	"5808C8AA-7E8F-42E0-85D2-E1E90434CFB3": "Windows LDM metadata",
	"E75CAF8F-F680-4CEE-AFA3-B001E56EFC2D": "Windows Storage Spaces",
	"37AFFC90-EF7D-4E96-91C3-2D7AE055B174": "IBM GPFS",
	"FE3A2A5D-4F32-41A7-B725-ACCC3285A309": "ChromeOS kernel",
	"3CB8E202-3B7E-47DD-8A3C-7FF2A13CFCEC": "ChromeOS root",
	"2E0A753D-9E48-43B0-8337-B15192CB1B5E": "ChromeOS reserved",
	"0657FD6D-A4AB-43C4-84E5-0933C84B4F4F": "Linux swap",
	"0FC63DAF-8483-4772-8E79-3D69D8477DE4": "Linux filesystem",
	"8DA63339-0007-60C0-C436-083AC8230908": "Linux reserved",
	"933AC7E1-2EB4-4F13-B844-0E14E2AEF915": "Linux /home",
	"3B8F8425-20E0-4F3B-907F-1A25A76F98E8": "Linux /srv",
	"7FFEC5C9-2D00-49B7-8941-3EA10A5586B7": "Linux dm-crypt",
	"CA7D7CCB-63ED-4C53-861C-1742536059CC": "Linux LUKS",
	"44479540-F297-41B2-9AF7-D131D5F0458A": "Linux x86",
	"4F68BCE3-E8CD-4DB1-96E7-FBCAF984B709": "Linux x86-64",
	"69DAD710-2CE4-4E3C-B16C-21A1D49ABED3": "Linux ARM32",
	"B921B045-1DF0-41C3-AF44-4C6F280D3FAE": "Linux ARM64",
	"993d8d3d-f80e-4225-855a-9daf8ed7ea97": "Linux IA-64",
	"D3BFE2DE-3DAF-11DF-BA40-E3A556D89593": "Intel Rapid Start",
	"E6D6D379-F507-44C2-A23C-238F2A3DF928": "Linux LVM",
	"734E5AFE-F61A-11E6-BC64-92361F002671": "Atari TOS",
	"516E7CB4-6ECF-11D6-8FF8-00022D09712B": "FreeBSD Disklabel",
	"83BD6B9D-7F41-11DC-BE0B-001560B84F0F": "FreeBSD boot",
	"516E7CB5-6ECF-11D6-8FF8-00022D09712B": "FreeBSD swap",
	"516E7CB6-6ECF-11D6-8FF8-00022D09712B": "FreeBSD UFS",
	"516E7CBA-6ECF-11D6-8FF8-00022D09712B": "FreeBSD ZFS",
	"516E7CB8-6ECF-11D6-8FF8-00022D09712B": "FreeBSD Vinum/RAID",
	"85D5E45A-237C-11E1-B4B3-E89A8F7FC3A7": "MidnightBSD data",
	"85D5E45E-237C-11E1-B4B3-E89A8F7FC3A7": "MidnightBSD boot",
	"85D5E45B-237C-11E1-B4B3-E89A8F7FC3A7": "MidnightBSD swap",
	"0394EF8B-237E-11E1-B4B3-E89A8F7FC3A7": "MidnightBSD UFS",
	"85D5E45D-237C-11E1-B4B3-E89A8F7FC3A7": "MidnightBSD ZFS",
	"85D5E45C-237C-11E1-B4B3-E89A8F7FC3A7": "MidnightBSD Vinum",
	"824CC7A0-36A8-11E3-890A-952519AD3F61": "OpenBSD data",
	"55465300-0000-11AA-AA11-00306543ECAC": "Apple UFS",
	"49F48D32-B10E-11DC-B99B-0019D1879648": "NetBSDNetBSD swap",
	"49F48D5A-B10E-11DC-B99B-0019D1879648": "NetBSD FFS",
	"49F48D82-B10E-11DC-B99B-0019D1879648": "NetBSD LFS",
	"2DB519C4-B10F-11DC-B99B-0019D1879648": "NetBSD concatenated",
	"2DB519EC-B10F-11DC-B99B-0019D1879648": "NetBSD encrypted",
	"49F48DAA-B10E-11DC-B99B-0019D1879648": "NetBSD RAID",
	"426F6F74-0000-11AA-AA11-00306543ECAC": "Apple boot",
	"48465300-0000-11AA-AA11-00306543ECAC": "Apple HFS/HFS+",
	"52414944-0000-11AA-AA11-00306543ECAC": "Apple RAID",
	"52414944-5F4F-11AA-AA11-00306543ECAC": "Apple RAID offline",
	"4C616265-6C00-11AA-AA11-00306543ECAC": "Apple Label",
	"5265636F-7665-11AA-AA11-00306543ECAC": "AppleTV Recovery",
	"53746F72-6167-11AA-AA11-00306543ECAC": "Apple Core Storage",
	"B6FA30DA-92D2-4A9A-96F1-871EC6486200": "Apple SoftRAID Status",
	"2E313465-19B9-463F-8126-8A7993773801": "Apple SoftRAID Scratch",
	"FA709C7E-65B1-4593-BFD5-E71D61DE9B02": "Apple SoftRAID Volume",
	"BBBA6DF5-F46F-4A89-8F59-8765B2727503": "Apple SoftRAID Cache",
	"7C3457EF-0000-11AA-AA11-00306543ECAC": "Apple APFS",
	"CEF5A9AD-73BC-4601-89F3-CDEEEEE321A1": "QNX6 Power-Safe",
	"0311FC50-01CA-4725-AD77-9ADBB20ACE98": "Acronis Secure Zone",
	"6A82CB45-1DD2-11B2-99A6-080020736631": "Solaris boot",
	"6A85CF4D-1DD2-11B2-99A6-080020736631": "Solaris root",
	"6A898CC3-1DD2-11B2-99A6-080020736631": "Solaris /usr",
	"6A87C46F-1DD2-11B2-99A6-080020736631": "Solaris swap",
	"6A8B642B-1DD2-11B2-99A6-080020736631": "Solaris backup",
	"6A8EF2E9-1DD2-11B2-99A6-080020736631": "Solaris /var",
	"6A90BA39-1DD2-11B2-99A6-080020736631": "Solaris /home",
	"6A9283A5-1DD2-11B2-99A6-080020736631": "Solaris alternate sector",
	"6A945A3B-1DD2-11B2-99A6-080020736631": "Solaris Reserved",
	"75894C1E-3AEB-11D3-B7C1-7B03A0000000": "HP-UX data",
	"E2A1E728-32E3-11D6-A682-7B03A0000000": "HP-UX service",
	"BC13C2FF-59E6-4262-A352-B275FD6F7172": "Freedesktop $BOOT",
	"42465331-3BA3-10F1-802A-4861696B7521": "Haiku BFS",
	"BFBFAFE7-A34F-448A-9A5B-6213EB736C22": "Lenovo system partition",
	"F4019732-066E-4E12-8273-346C5641494F": "Sony system partition",
	"C12A7328-F81F-11D2-BA4B-00A0C93EC93B": "EFI System (ESP)",
	"024DEE41-33E7-11D3-9D69-0008C781F39F": "MBR partition scheme",
	"21686148-6449-6E6F-744E-656564454649": "BIOS boot partition",
	"4FBD7E29-9D25-41B8-AFD0-062C0CEFF05D": "Ceph OSD",
	"4FBD7E29-9D25-41B8-AFD0-5EC00CEFF05D": "Ceph dm-crypt OSD",
	"45B0969E-9B03-4F30-B4C6-B4B80CEFF106": "Ceph journal",
	"45B0969E-9B03-4F30-B4C6-5EC00CEFF106": "Ceph dm-crypt journal",
	"89C57F98-2FE5-4DC0-89C1-F3AD0CEFF2BE": "Ceph disk in creation",
	"89C57F98-2FE5-4DC0-89C1-5EC00CEFF2BE": "Ceph dm-crypt disk in creation",
	"AA31E02A-400F-11DB-9590-000C2911D1B8": "VMware ESX VMFS",
	"9198EFFC-31C0-11DB-8F78-000C2911D1B8": "VMware reserved",
	"9D275380-40AD-11DB-BF97-000C2911D1B8": "VMware kcore crash protection",
	"A19D880F-05FC-4D3B-A006-743F0F84911E": "Linux RAID",
}
