package segword

import (
	"fmt"
	"testing"
)

func TestSegWords(t *testing.T) {
	contents := []string{
		`1/2/3`,
		`宇宙（Universe）在物理意义上被定义为所有的空间和时间（统称为时空）及其内涵，包括各种形式的所有能量，比如电磁辐射、普通物质、暗物质、暗能量等，其中普通物质包括行星、卫星、恒星、星系、星系团和星系间物质等。宇宙还包括影响物质和能量的物理定律，如守恒定律、经典力学、相对论等。 [1] `,
		`大爆炸理论是关于宇宙演化的现代宇宙学描述。根据这一理论的估计，空间和时间在137.99±0.21亿年前的大爆炸后一同出现，随着宇宙膨胀，最初存在的能量和物质变得不那么密集。最初的加速膨胀被称为暴胀时期，之后已知的四个基本力分离。宇宙逐渐冷却并继续膨胀，允许第一个亚原子粒子和简单的原子形成。暗物质逐渐聚集，在引力作用下形成泡沫一样的结构，大尺度纤维状结构和宇宙空洞。巨大的氢氦分子云逐渐被吸引到暗物质最密集的地方，形成了第一批星系、恒星、行星以及所有的一切。空间本身在不断膨胀，因此当前可以看见距离地球465亿光年的天体，因为这些光在138亿年前产生的时候距离地球比当前更近。`,
	}

	t.Run("all", func(t *testing.T) {
		sw := NewSegWords(NewExportAll())
		err := sw.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := sw.Do(contents)
		rets := token.ToToken()
		if len(rets) <= 0 {
			t.Error("tag error")
		}

		fmt.Println(rets)
	})

	t.Run("line", func(t *testing.T) {
		sw := NewSegWords(NewExportLine())
		err := sw.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := sw.Do(contents)
		rets := token.ToToken()
		if len(rets) <= 0 {
			t.Error("tag error")
		}

		fmt.Println(rets)
	})
}
