package eip

func swuMapForG1(u fe, fq *fq, params *g1SWUParams) (fe, fe, bool) {
	var tv [4]fe
	for i := 0; i < 4; i++ {
		tv[i] = fq.new()
	}

	// 1.  tv1 = Z * u^2
	fq.square(tv[0], u)
	fq.mul(tv[0], tv[0], params.z)

	// 2.  tv2 = tv1^2
	fq.square(tv[1], tv[0])

	// 3.   x1 = tv1 + tv2
	x1 := fq.new()
	fq.add(x1, tv[0], tv[1])

	// 4.   x1 = inv0(x1)
	fq.inverse(x1, x1)

	// 5.   e1 = x1 == 0
	e1 := fq.isZero(x1)

	// 6.   x1 = x1 + 1
	fq.add(x1, x1, fq.one)

	// 7.   x1 = CMOV(x1, c2, e1)    # If (tv1 + tv2) == 0, set x1 = -1 / Z
	if e1 {
		fq.copy(x1, params.zInv)
	}

	// 8.   x1 = x1 * c1      # x1 = (-B / A) * (1 + (1 / (Z^2 * u^4 + Z * u^2)))
	fq.mul(x1, x1, params.minusBOverA)

	// 9.  gx1 = x1^2
	gx1 := fq.new()
	fq.square(gx1, x1)

	// 10. gx1 = gx1 + A
	fq.add(gx1, gx1, params.a) // TODO: a is zero we can ommit

	// 11. gx1 = gx1 * x1
	fq.mul(gx1, gx1, x1)

	// 12. gx1 = gx1 + B             # gx1 = g(x1) = x1^3 + A * x1 + B
	fq.add(gx1, gx1, params.b)

	// 13.  x2 = tv1 * x1            # x2 = Z * u^2 * x1
	x2 := fq.new()
	fq.mul(x2, tv[0], x1)

	//  // 14. tv2 = tv1 * tv2
	fq.mul(tv[1], tv[0], tv[1])

	//  // 15. gx2 = gx1 * tv2           # gx2 = (Z * u^2)^3 * gx1
	gx2 := fq.new()
	fq.mul(gx2, gx1, tv[1])

	// 16.  e2 = is_square(gx1)
	_e2 := legendreSymbolFq(fq, gx1)
	e2 := _e2 == 0 || _e2 == 1
	// 17.   x = CMOV(x2, x1, e2)    # If is_square(gx1), x = x1, else x = x2
	x := fq.new()
	if e2 {
		fq.copy(x, x1)
	} else {
		fq.copy(x, x2)
	}

	// 18.  y2 = CMOV(gx2, gx1, e2)  # If is_square(gx1), y2 = gx1, else y2 = gx2
	y2 := fq.new()
	if e2 {
		fq.copy(y2, gx1)
	} else {
		fq.copy(y2, gx2)
	}

	// 19.   y = sqrt(y2)
	y := fq.new()
	if hasSquareRoot := fq.sqrt(y, y2); !hasSquareRoot {
		return nil, nil, false
	}
	// 20.  e3 = sgn0(u) == sgn0(y)  # Fix sign of y
	uSign := fq.sign(u)
	ySign := fq.sign(y)

	if ((uSign == 1 && ySign == -1) || (uSign == -1 && ySign == 1)) || ((uSign == 0 && ySign == -1) || (uSign == -1 && ySign == 0)) {
		fq.neg(y, y)
	}
	return x, y, true
}

func swuMapForG2(u *fe2, fq2 *fq2, params *g2SWUParams) (*fe2, *fe2, bool) {
	var tv [4]*fe2
	for i := 0; i < 4; i++ {
		tv[i] = fq2.new()
	}

	// 1.  tv1 = Z * u^2
	fq2.square(tv[0], u)
	fq2.mul(tv[0], tv[0], params.z)

	// 2.  tv2 = tv1^2
	fq2.square(tv[1], tv[0])

	// 3.   x1 = tv1 + tv2
	x1 := fq2.new()
	fq2.add(x1, tv[0], tv[1])

	// 4.   x1 = inv0(x1)
	fq2.inverse(x1, x1)

	// 5.   e1 = x1 == 0
	e1 := fq2.isZero(x1)

	// 6.   x1 = x1 + 1
	fq2.add(x1, x1, fq2.one())

	// 7.   x1 = CMOV(x1, c2, e1)    # If (tv1 + tv2) == 0, set x1 = -1 / Z
	if e1 {
		fq2.copy(x1, params.zInv)
	}

	// 8.   x1 = x1 * c1      # x1 = (-B / A) * (1 + (1 / (Z^2 * u^4 + Z * u^2)))
	fq2.mul(x1, x1, params.minusBOverA)

	// 9.  gx1 = x1^2
	gx1 := fq2.new()
	fq2.square(gx1, x1)
	// 10. gx1 = gx1 + A
	fq2.add(gx1, gx1, params.a) // TODO: a is zero we can ommit

	// 11. gx1 = gx1 * x1
	fq2.mul(gx1, gx1, x1)

	// 12. gx1 = gx1 + B             # gx1 = g(x1) = x1^3 + A * x1 + B
	fq2.add(gx1, gx1, params.b)

	// 13.  x2 = tv1 * x1            # x2 = Z * u^2 * x1
	x2 := fq2.new()
	fq2.mul(x2, tv[0], x1)

	// 14. tv2 = tv1 * tv2
	fq2.mul(tv[1], tv[0], tv[1])

	// 15. gx2 = gx1 * tv2           # gx2 = (Z * u^2)^3 * gx1
	gx2 := fq2.new()
	fq2.mul(gx2, gx1, tv[1])

	// 16.  e2 = is_square(gx1)
	_e2 := legendreSymbolFq2(fq2, gx1)
	e2 := _e2 == 0 || _e2 == 1
	// 17.   x = CMOV(x2, x1, e2)    # If is_square(gx1), x = x1, else x = x2
	x := fq2.new()
	if e2 {
		fq2.copy(x, x1)
	} else {
		fq2.copy(x, x2)
	}

	// 18.  y2 = CMOV(gx2, gx1, e2)  # If is_square(gx1), y2 = gx1, else y2 = gx2
	y2 := fq2.new()
	if e2 {
		fq2.copy(y2, gx1)
	} else {
		fq2.copy(y2, gx2)
	}

	// 19.   y = sqrt(y2)
	y := fq2.new()
	if hasSquareRoot := fq2.sqrt(y, y2); !hasSquareRoot {
		return nil, nil, false
	}
	// 20.  e3 = sgn0(u) == sgn0(y)  # Fix sign of y
	uSign := fq2.sign(u)
	ySign := fq2.sign(y)

	if ((uSign == 1 && ySign == -1) || (uSign == -1 && ySign == 1)) || ((uSign == 0 && ySign == -1) || (uSign == -1 && ySign == 0)) {
		fq2.neg(y, y)
	}
	return x, y, true
}

func legendreSymbolFq(fq *fq, elem fe) int {
	if fq.isZero(elem) {
		return 0
	} else if fq.isNonResidue(elem, 2) {
		return -1
	} else {
		return 1
	}
}

func legendreSymbolFq2(fq2 *fq2, elem *fe2) int {
	if fq2.isZero(elem) {
		return 0
	} else if fq2.isNonResidue(elem, 2) {
		return -1
	} else {
		return 1
	}
}

func applyIsogenyMapForG1(fq *fq, x, y fe, params *g1IsogenyParams) (fe, fe) {
	degree := 15
	xNum := params.k1[degree]
	xDen := params.k2[degree]
	yNum := params.k3[degree]
	yDen := params.k4[degree]

	for i := degree - 1; i >= 0; i-- {
		fq.mul(xNum, xNum, x)
		fq.mul(xDen, xDen, x)
		fq.mul(yNum, yNum, x)
		fq.mul(yDen, yDen, x)

		fq.add(xNum, xNum, params.k1[i])
		fq.add(xDen, xDen, params.k2[i])
		fq.add(yNum, yNum, params.k3[i])
		fq.add(yDen, yDen, params.k4[i])
	}

	fq.inverse(xDen, xDen)
	fq.inverse(yDen, yDen)

	fq.mul(xNum, xNum, xDen)
	fq.mul(yNum, yNum, yDen)

	fq.mul(yNum, yNum, y)

	return xNum, yNum
}

func applyIsogenyMapForG2(fq2 *fq2, x, y *fe2, params *g2IsogenyParams) (*fe2, *fe2) {
	degree := 3
	xNum := params.k1[degree]
	xDen := params.k2[degree]
	yNum := params.k3[degree]
	yDen := params.k4[degree]

	for i := degree - 1; i >= 0; i-- {
		fq2.mul(xNum, xNum, x)
		fq2.mul(xDen, xDen, x)
		fq2.mul(yNum, yNum, x)
		fq2.mul(yDen, yDen, x)

		fq2.add(xNum, xNum, params.k1[i])
		fq2.add(xDen, xDen, params.k2[i])
		fq2.add(yNum, yNum, params.k3[i])
		fq2.add(yDen, yDen, params.k4[i])
	}

	fq2.inverse(xDen, xDen)
	fq2.inverse(yDen, yDen)

	fq2.mul(xNum, xNum, xDen)
	fq2.mul(yNum, yNum, yDen)

	fq2.mul(yNum, yNum, y)

	return xNum, yNum
}

type g1SWUParams struct {
	z           fe
	zInv        fe
	a           fe
	b           fe
	minusBOverA fe
}

type g2SWUParams struct {
	z           *fe2
	zInv        *fe2
	a           *fe2
	b           *fe2
	minusBOverA *fe2
}

func computeSWUParamsForG1(fq *fq) *g1SWUParams {
	z, _ := fq.fromString("0x0b")
	a, _ := fq.fromString("0x00144698a3b8e9433d693a02c96d4982b0ea985383ee66a8d8e8981aefd881ac98936f8da0e0f97f5cf428082d584c1d")
	b, _ := fq.fromString("0x12e2908d11688030018b12e8753eee3b2016c1f0f24f4070a0b9c14fcef35ef55a23215a316ceaa5d1cc48e98e172be0")
	// z^-1
	zInv, _ := fq.fromString("0x025d302c90dd14f6c102839c34a9c9e509221e235bf4d328ac4a41b18aca44ec02c9d1743eaa8ba2e25cfffffffff83e")
	// -b/a
	minusBOverA, _ := fq.fromString("0x0793154fd85631d966ef2470460c78f6a928ad9f5bdbfac21df39753aa278ba751bdfcf95a84188e29d670675e4c9c7c")

	return &g1SWUParams{
		z,
		zInv,
		a,
		b,
		minusBOverA,
	}

}

func computeSWUParamsForG2(fq2 *fq2) *g2SWUParams {
	z, a, b, zInv, minusBOverA := fq2.new(), fq2.new(), fq2.new(), fq2.new(), fq2.new()

	z[0], _ = fq2.f.fromString("0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaa9")
	z[1], _ = fq2.f.fromString("0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaaa")
	a[0], _ = fq2.f.fromString("0x00")
	a[1], _ = fq2.f.fromString("0xf0")
	b[0], _ = fq2.f.fromString("0x03f4")
	b[1], _ = fq2.f.fromString("0x03f4")
	// z^-1
	zInv[0], _ = fq2.f.fromString("0x053369fba51994854238bb2473dbef5e474b0f1a971a9d597b09c3b9caf0313a6c88cccc89dd99998b99666666665555")
	zInv[1], _ = fq2.f.fromString("0x0a66d3f74a33290a84717648e7b7debc8e961e352e353ab2f613877395e06274d911999913bb33331732ccccccccaaab")
	// -b/a
	minusBOverA[0], _ = fq2.f.fromString("0x083c12791abdd5d2fe2f284f0cc6e5aa9b8c2d3f6f3f792302cf75e62bfc4df1d6834443da498888725d8cccccccb1c3")
	minusBOverA[1], _ = fq2.f.fromString("0x11c4ff711ec210c74cec7f673684c72cc8eb1e458445999c64615cbacab4a8324828bbbad70a777747a173333332f8e8")

	return &g2SWUParams{
		z,
		zInv,
		a,
		b,
		minusBOverA,
	}
}

type g1IsogenyParams struct {
	k1 [16]fe
	k2 [16]fe
	k3 [16]fe
	k4 [16]fe
}

type g2IsogenyParams struct {
	k1 [4]*fe2
	k2 [4]*fe2
	k3 [4]*fe2
	k4 [4]*fe2
}

func prepareIsogenyParamsForG1(fq *fq) *g1IsogenyParams {
	var k1 [16]fe
	k1Constants := [16]string{
		"0x11a05f2b1e833340b809101dd99815856b303e88a2d7005ff2627b56cdb4e2c85610c2d5f2e62d6eaeac1662734649b7",
		"0x17294ed3e943ab2f0588bab22147a81c7c17e75b2f6a8417f565e33c70d1e86b4838f2a6f318c356e834eef1b3cb83bb",
		"0x0d54005db97678ec1d1048c5d10a9a1bce032473295983e56878e501ec68e25c958c3e3d2a09729fe0179f9dac9edcb0",
		"0x1778e7166fcc6db74e0609d307e55412d7f5e4656a8dbf25f1b33289f1b330835336e25ce3107193c5b388641d9b6861",
		"0x0e99726a3199f4436642b4b3e4118e5499db995a1257fb3f086eeb65982fac18985a286f301e77c451154ce9ac8895d9",
		"0x1630c3250d7313ff01d1201bf7a74ab5db3cb17dd952799b9ed3ab9097e68f90a0870d2dcae73d19cd13c1c66f652983",
		"0x0d6ed6553fe44d296a3726c38ae652bfb11586264f0f8ce19008e218f9c86b2a8da25128c1052ecaddd7f225a139ed84",
		"0x17b81e7701abdbe2e8743884d1117e53356de5ab275b4db1a682c62ef0f2753339b7c8f8c8f475af9ccb5618e3f0c88e",
		"0x080d3cf1f9a78fc47b90b33563be990dc43b756ce79f5574a2c596c928c5d1de4fa295f296b74e956d71986a8497e317",
		"0x169b1f8e1bcfa7c42e0c37515d138f22dd2ecb803a0c5c99676314baf4bb1b7fa3190b2edc0327797f241067be390c9e",
		"0x10321da079ce07e272d8ec09d2565b0dfa7dccdde6787f96d50af36003b14866f69b771f8c285decca67df3f1605fb7b",
		"0x06e08c248e260e70bd1e962381edee3d31d79d7e22c837bc23c0bf1bc24c6b68c24b1b80b64d391fa9c8ba2e8ba2d229",
		"0x00",
		"0x00",
		"0x00",
		"0x00",
	}

	var k2 [16]fe
	k2Constants := [16]string{
		"0x08ca8d548cff19ae18b2e62f4bd3fa6f01d5ef4ba35b48ba9c9588617fc8ac62b558d681be343df8993cf9fa40d21b1c",
		"0x12561a5deb559c4348b4711298e536367041e8ca0cf0800c0126c2588c48bf5713daa8846cb026e9e5c8276ec82b3bff",
		"0x0b2962fe57a3225e8137e629bff2991f6f89416f5a718cd1fca64e00b11aceacd6a3d0967c94fedcfcc239ba5cb83e19",
		"0x03425581a58ae2fec83aafef7c40eb545b08243f16b1655154cca8abc28d6fd04976d5243eecf5c4130de8938dc62cd8",
		"0x13a8e162022914a80a6f1d5f43e7a07dffdfc759a12062bb8d6b44e833b306da9bd29ba81f35781d539d395b3532a21e",
		"0x0e7355f8e4e667b955390f7f0506c6e9395735e9ce9cad4d0a43bcef24b8982f7400d24bc4228f11c02df9a29f6304a5",
		"0x0772caacf16936190f3e0c63e0596721570f5799af53a1894e2e073062aede9cea73b3538f0de06cec2574496ee84a3a",
		"0x14a7ac2a9d64a8b230b3f5b074cf01996e7f63c21bca68a81996e1cdf9822c580fa5b9489d11e2d311f7d99bbdcc5a5e",
		"0x0a10ecf6ada54f825e920b3dafc7a3cce07f8d1d7161366b74100da67f39883503826692abba43704776ec3a79a1d641",
		"0x095fc13ab9e92ad4476d6e3eb3a56680f682b4ee96f7d03776df533978f31c1593174e4b4b7865002d6384d168ecdd0a",
		"0x01",
		"0x00",
		"0x00",
		"0x00",
		"0x00",
		"0x00",
	}

	var k3 [16]fe
	k3Constants := []string{
		"0x090d97c81ba24ee0259d1f094980dcfa11ad138e48a869522b52af6c956543d3cd0c7aee9b3ba3c2be9845719707bb33",
		"0x134996a104ee5811d51036d776fb46831223e96c254f383d0f906343eb67ad34d6c56711962fa8bfe097e75a2e41c696",
		"0x00cc786baa966e66f4a384c86a3b49942552e2d658a31ce2c344be4b91400da7d26d521628b00523b8dfe240c72de1f6",
		"0x01f86376e8981c217898751ad8746757d42aa7b90eeb791c09e4a3ec03251cf9de405aba9ec61deca6355c77b0e5f4cb",
		"0x08cc03fdefe0ff135caf4fe2a21529c4195536fbe3ce50b879833fd221351adc2ee7f8dc099040a841b6daecf2e8fedb",
		"0x16603fca40634b6a2211e11db8f0a6a074a7d0d4afadb7bd76505c3d3ad5544e203f6326c95a807299b23ab13633a5f0",
		"0x04ab0b9bcfac1bbcb2c977d027796b3ce75bb8ca2be184cb5231413c4d634f3747a87ac2460f415ec961f8855fe9d6f2",
		"0x0987c8d5333ab86fde9926bd2ca6c674170a05bfe3bdd81ffd038da6c26c842642f64550fedfe935a15e4ca31870fb29",
		"0x09fc4018bd96684be88c9e221e4da1bb8f3abd16679dc26c1e8b6e6a1f20cabe69d65201c78607a360370e577bdba587",
		"0x0e1bba7a1186bdb5223abde7ada14a23c42a0ca7915af6fe06985e7ed1e4d43b9b3f7055dd4eba6f2bafaaebca731c30",
		"0x19713e47937cd1be0dfd0b8f1d43fb93cd2fcbcb6caf493fd1183e416389e61031bf3a5cce3fbafce813711ad011c132",
		"0x18b46a908f36f6deb918c143fed2edcc523559b8aaf0c2462e6bfe7f911f643249d9cdf41b44d606ce07c8a4d0074d8e",
		"0x0b182cac101b9399d155096004f53f447aa7b12a3426b08ec02710e807b4633f06c851c1919211f20d4c04f00b971ef8",
		"0x0245a394ad1eca9b72fc00ae7be315dc757b3b080d4c158013e6632d3c40659cc6cf90ad1c232a6442d9d3f5db980133",
		"0x05c129645e44cf1102a159f748c4a3fc5e673d81d7e86568d9ab0f5d396a7ce46ba1049b6579afb7866b1e715475224b",
		"0x15e6be4e990f03ce4ea50b3b42df2eb5cb181d8f84965a3957add4fa95af01b2b665027efec01c7704b456be69c8b604",
	}

	var k4 [16]fe
	k4Constants := []string{
		"0x16112c4c3a9c98b252181140fad0eae9601a6de578980be6eec3232b5be72e7a07f3688ef60c206d01479253b03663c1",
		"0x1962d75c2381201e1a0cbd6c43c348b885c84ff731c4d59ca4a10356f453e01f78a4260763529e3532f6102c2e49a03d",
		"0x058df3306640da276faaae7d6e8eb15778c4855551ae7f310c35a5dd279cd2eca6757cd636f96f891e2538b53dbf67f2",
		"0x16b7d288798e5395f20d23bf89edb4d1d115c5dbddbcd30e123da489e726af41727364f2c28297ada8d26d98445f5416",
		"0x0be0e079545f43e4b00cc912f8228ddcc6d19c9f0f69bbb0542eda0fc9dec916a20b15dc0fd2ededda39142311a5001d",
		"0x08d9e5297186db2d9fb266eaac783182b70152c65550d881c5ecd87b6f0f5a6449f38db9dfa9cce202c6477faaf9b7ac",
		"0x166007c08a99db2fc3ba8734ace9824b5eecfdfa8d0cf8ef5dd365bc400a0051d5fa9c01a58b1fb93d1a1399126a775c",
		"0x16a3ef08be3ea7ea03bcddfabba6ff6ee5a4375efa1f4fd7feb34fd206357132b920f5b00801dee460ee415a15812ed9",
		"0x1866c8ed336c61231a1be54fd1d74cc4f9fb0ce4c6af5920abc5750c4bf39b4852cfe2f7bb9248836b233d9d55535d4a",
		"0x167a55cda70a6e1cea820597d94a84903216f763e13d87bb5308592e7ea7d4fbc7385ea3d529b35e346ef48bb8913f55",
		"0x04d2f259eea405bd48f010a01ad2911d9c6dd039bb61a6290e591b36e636a5c871a5c29f4f83060400f8b49cba8f6aa8",
		"0x0accbb67481d033ff5852c1e48c50c477f94ff8aefce42d28c0f9a88cea7913516f968986f7ebbea9684b529e2561092",
		"0x0ad6b9514c767fe3c3613144b45f1496543346d98adf02267d5ceef9a00d9b8693000763e3b90ac11e99b138573345cc",
		"0x02660400eb2e4f3b628bdd0d53cd76f2bf565b94e72927c1cb748df27942480e420517bd8714cc80d1fadc1326ed06f7",
		"0x0e0fa1d816ddc03e6b24255e0d7819c171c40f65e273b853324efcd6356caa205ca2f570f13497804415473a1d634b8f",
		"0x01",
	}

	for i := 0; i < len(k1Constants); i++ {
		k1[i], _ = fq.fromString(k1Constants[i])
		k2[i], _ = fq.fromString(k2Constants[i])
		k3[i], _ = fq.fromString(k3Constants[i])
		k4[i], _ = fq.fromString(k4Constants[i])
	}

	return &g1IsogenyParams{
		k1,
		k2,
		k3,
		k4,
	}
}

func prepareIsogenyParamsForG2(fq2 *fq2) *g2IsogenyParams {
	var k1, k2, k3, k4 [4]*fe2
	for i := 0; i < 4; i++ {
		k1[i], k2[i], k3[i], k4[i] = fq2.new(), fq2.new(), fq2.new(), fq2.new()
	}

	k1[0][0], _ = fq2.f.fromString("0x05c759507e8e333ebb5b7a9a47d7ed8532c52d39fd3a042a88b58423c50ae15d5c2638e343d9c71c6238aaaaaaaa97d6")
	k1[0][1], _ = fq2.f.fromString("0x05c759507e8e333ebb5b7a9a47d7ed8532c52d39fd3a042a88b58423c50ae15d5c2638e343d9c71c6238aaaaaaaa97d6")
	k1[1][0], _ = fq2.f.fromString("0x00")
	k1[1][1], _ = fq2.f.fromString("0x11560bf17baa99bc32126fced787c88f984f87adf7ae0c7f9a208c6b4f20a4181472aaa9cb8d555526a9ffffffffc71a")
	k1[2][0], _ = fq2.f.fromString("0x11560bf17baa99bc32126fced787c88f984f87adf7ae0c7f9a208c6b4f20a4181472aaa9cb8d555526a9ffffffffc71e")
	k1[2][1], _ = fq2.f.fromString("0x08ab05f8bdd54cde190937e76bc3e447cc27c3d6fbd7063fcd104635a790520c0a395554e5c6aaaa9354ffffffffe38d")
	k1[3][0], _ = fq2.f.fromString("0x171d6541fa38ccfaed6dea691f5fb614cb14b4e7f4e810aa22d6108f142b85757098e38d0f671c7188e2aaaaaaaa5ed1")
	k1[3][1], _ = fq2.f.fromString("0x00")

	k2[0][0], _ = fq2.f.fromString("0x00")
	k2[0][1], _ = fq2.f.fromString("0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaa63")
	k2[1][0], _ = fq2.f.fromString("0x0c")
	k2[1][1], _ = fq2.f.fromString("0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaa9f")
	k2[2][0], _ = fq2.f.fromString("0x01")
	k2[2][1], _ = fq2.f.fromString("0x00")

	k3[0][0], _ = fq2.f.fromString("0x1530477c7ab4113b59a4c18b076d11930f7da5d4a07f649bf54439d87d27e500fc8c25ebf8c92f6812cfc71c71c6d706")
	k3[0][1], _ = fq2.f.fromString("0x1530477c7ab4113b59a4c18b076d11930f7da5d4a07f649bf54439d87d27e500fc8c25ebf8c92f6812cfc71c71c6d706")
	k3[1][0], _ = fq2.f.fromString("0x00")
	k3[1][1], _ = fq2.f.fromString("0x05c759507e8e333ebb5b7a9a47d7ed8532c52d39fd3a042a88b58423c50ae15d5c2638e343d9c71c6238aaaaaaaa97be")
	k3[2][0], _ = fq2.f.fromString("0x11560bf17baa99bc32126fced787c88f984f87adf7ae0c7f9a208c6b4f20a4181472aaa9cb8d555526a9ffffffffc71c")
	k3[2][1], _ = fq2.f.fromString("0x08ab05f8bdd54cde190937e76bc3e447cc27c3d6fbd7063fcd104635a790520c0a395554e5c6aaaa9354ffffffffe38f")
	k3[3][0], _ = fq2.f.fromString("0x124c9ad43b6cf79bfbf7043de3811ad0761b0f37a1e26286b0e977c69aa274524e79097a56dc4bd9e1b371c71c718b10")
	k3[3][1], _ = fq2.f.fromString("0x00")

	k4[0][0], _ = fq2.f.fromString("0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffa8fb")
	k4[0][1], _ = fq2.f.fromString("0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffa8fb")
	k4[1][0], _ = fq2.f.fromString("0x00")
	k4[1][1], _ = fq2.f.fromString("0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffa9d3")
	k4[2][0], _ = fq2.f.fromString("0x12")
	k4[2][1], _ = fq2.f.fromString("0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaa99")
	k4[3][0], _ = fq2.f.fromString("0x01")
	k4[3][1], _ = fq2.f.fromString("0x00")

	return &g2IsogenyParams{
		k1,
		k2,
		k3,
		k4,
	}

}
