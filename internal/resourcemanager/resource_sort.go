package resourcemanager

import "github.com/amidgo/cloud-resources/internal/model/resourcemodel"

// сортировка ресурсов по возрастанию CPU
type SortResourcesByCPUAsc []*resourcemodel.Resource

func (a SortResourcesByCPUAsc) Len() int           { return len(a) }
func (a SortResourcesByCPUAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortResourcesByCPUAsc) Less(i, j int) bool { return a[i].CPU < a[j].CPU }

// сортировка ресурсов по убыванию CPU
type SortResourcesByCPUDesc []*resourcemodel.Resource

func (a SortResourcesByCPUDesc) Len() int           { return len(a) }
func (a SortResourcesByCPUDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortResourcesByCPUDesc) Less(i, j int) bool { return a[i].CPU > a[j].CPU }

// сортировка ресурсов по возрастанию RAM
type SortResourcesByRAMAsc []*resourcemodel.Resource

func (a SortResourcesByRAMAsc) Len() int           { return len(a) }
func (a SortResourcesByRAMAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortResourcesByRAMAsc) Less(i, j int) bool { return a[i].RAM < a[j].RAM }

type SortResourcesByRAMDesc []*resourcemodel.Resource

func (a SortResourcesByRAMDesc) Len() int           { return len(a) }
func (a SortResourcesByRAMDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortResourcesByRAMDesc) Less(i, j int) bool { return a[i].RAM > a[j].RAM }
