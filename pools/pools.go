package pools

// Eth1 Addresses used to identify the pools in the deposit contract
var PoolsAddresses = map[string][]string{
	"kraken": {
		"0xa40dfee99e1c85dc97fdc594b16a460717838703",
		"0x631c2d8d0d7a80824e602a79800a98d93e909a9e",
	},
	"binance": {
		"0xBdD75A97c29294FF805FB2fEe65aBd99492b32A8",
		"0x50fF765A993400CD62B61Cfa4Bb33B1dDF694eC7",
		"0xf211Dbb151048f65895d99A446c5268198Af73D2",
		"0x071A5f87Ca7cbDd86Ba03ab11f68C8fA2A542B91",
		"0xfDE01891bC1DdA13Ad2B6027709777066290FD72",
		"0xf3Dca7dc9265D92b17A054Cf43fE6e02c571553b",
		"0x339c367c63d60aF280FB727140B7469675D478bc",
		"0xF4fEae08C1Fa864B64024238E33Bfb4A3Ea7741d",
		"0x996793790D726072273Ba5eEf1E15a032e847e2B",
		"0x8ebcf6Dbe031f6ce1f3FDe97747f1B4Ad6ab88B2",
		"0x32895327eEAbC019Bb3363A1cC14a078A3bb36fA",
		"0xA2ABe03F3A906dc11e05e90489946E5844374708",
		"0xe572d341Fdb292F0bf8964F3Ff9b0c2b9498f1C2",
		"0xb53f8c05e52e3c49283528aa24f3494375325931",
		"0x2f47a1c2db4a3b78cda44eade915c3b19107ddcc",
		"0x82361cb56a94c803593fde178bab2f345178e901",
		"0x0b75ef3fc3394548efed0d784afa32de81ad4923",
		"0xeae401546b85ae9b8bf4b81b4fb5c4337f079c09",
		"0x73483496eb9317ce4559f707c06b9377627a61b0",
		"0x32c96d17d81615789357160b41da2ef8b712eba8",
		"0xd1366d60d65bdfc92e5a5925fc4698a22e04e8c2",
		"0x90ba5b3a3e9353f713bb873a44f2396429c1dd27",
		"0xf17aced3c7a8daa29ebb90db8d1b6efd8c364a18",
	},
	"lido": {
		"0x00444797ba158a7bdb8302e72da98dcbccef0fbc",
		"0x073adf97f6de257d76e67f7c2fe57ac9843cca25",
		"0x0ac7e9af32422ac5968622585822e4d89ef51343",
		"0x26bc5cf2269cac024b106856ad37a1439cadd731",
		"0x41e37ddde5c80ee80cb9c468c0ed342078b78f2e",
		"0x41eff4e547090fbd310099fe13d18069c65dfe61",
		"0x4f1512d88d00962f58b6d543ff1803cc683d39f5",
		"0x55bc991b2edf3ddb4c520b222be4f378418ff0fa",
		"0x56d1eb6d70501b5dc5f1127ebaa80ee9a15a114c",
		"0x5fcf6e9c7e97986a14a79a8c7a69b928a7812dcd",
		"0x6352f8c749954c9df198cf72976e48994a77cce2",
		"0x68575571e75d2cfa4222e0f8e7053f056eb91d6c",
		"0x6e34e47df7026e0ace9457f930f1cfada6f547c4",
		"0x99b2c5d50086b02f83e791633c5660fbb8344653",
		"0xa2cd5b9d50d19fdfd2f37cbae0e880f9ce327837",
		"0xa76a7d0d06754e4fc4941519d1f9d56fd9f8d53b",
		"0xadf76760f1d6b984e39c27c87d5f9661cefc5a21",
		"0xb049e2336cacf0ba1f735fb8303c8aab227b90d2",
		"0xb96f132ad2e293a47619ec34b9d090cfe9735820",
		"0xbbe7fb58301d5abb45419342a7dd5c6cd68e2ff3",
		"0xc76050039bb278a9ece4ade1670c5ed91804cbc1",
		"0xc8381ca290c198f5ab739a1841ce8aedb0b330d5",
		"0xcc85da708a4ee036f2e1cf4b75dcadf49a3382cf",
		"0xd9f84209d3caa6b0a2ede8fbbe9fe4241e64cbb8",
		"0xdd19274b614b5ecacf493bc43c380ef6b8dfb56c",
		"0xe19fc582dd93fa876cf4061eb5456f310144f57b",
		"0xfbec500765ea318658de235bf8219eeb1c8fa540",
		"0xfe56a0dbdad44dd14e4d560632cc842c8a13642b",
		"0xfff8a72c72e0d5e08e85be05868990e8e4eef2da",
		"0xf82ac5937a20dc862f9bc0668779031e06000f17",
	},
	"huobi": {
		"0x194bd70b59491ce1310ea0bceabdb6c23ac9d5b2",
		"0xb73f4d4e99f65ec4b16b684e44f81aeca5ba2b7c",
	},
	"bitcoinsuisse": {
		"0xc2288b408dc872a1546f13e6ebfa9c94998316a2",
		"0xdd9663bd979f1ab1bada85e1bc7d7f13cafe71f8",
		"0x622de9bb9ff8907414785a633097db438f9a2d86",
		"0x3837ea2279b8e5c260a78f5f4181b783bbe76a8b",
		"0xec70e3c8afe212039c3f6a2df1c798003bf7cfe9",
		"0x2a7077399b3e90f5392d55a1dc7046ad8d152348",
	},
	"stakedus": {
		"0x54a9411c60a38054858ef227c3201e579bac5ff9",
		"0x13767b8816e6e03325621189311429d39c8325d1",
		"0xb770050098ea9870e030a1b460761d449e99095a",
		"0xc6b122b1933b0135a490b45e20f23803e78a73d2",
		"0xaa14fa2bc16b5ed1224625d746c2c83ee7f293c7",
		"0x6d6d0c7ff7cb9bbfa42f3ac4e7804a398a9f7fac",
		"0x2f9fd66e8894cec539c963b62e01baf4e7db018a",
		"0x7f78c085441419f163a9a25522052240040366ef",
		"0x8083c8cc154c22d50012a4651adabb1ac2d15f38",
		"0x05a1ff0a32bc24265bcb39499d0c5d9a6cb2011c",
		"0x03e46a1c771f3cf016462070bba3749b0870e426",
		"0xa0cd2aa0dad8c5a1b8c9f0549a9d042b7f9a9578",
		"0xd396e4f302c55c990cf67141ac4ee9d4fb88a67d",
		"0x16020b32e2f486416cd979afcaa894ce0d0d9fd7",
		"0x8657f4a1390792e72ba9ea5a868ae3801cc692ae",
		"0x95ad18082718659a207fa82146fed9217ae71588",
		"0x0d4ef0dbc65488c2e47e3700795c9c9f5bf0bf1b",
		"0x11f856a61e87a9b317fc20c8da9af13839c8658a",
		"0x307da664088fbb1728e8827bfc1a4459cc6a452c",
		"0x180fe42a404f406b943068dcdabb9c541dd6c8d6",
		"0xb342c7fa3d88a8e622f44ea335d3e97f91ef4e17",
		"0x2fba8be27cb36f30fdc258808ec6f32fc22222b7",
		"0xc74cb4dc7b27bbe018c7ce7a86b93ce859b2461a",
		"0x9f5a7c32be14413349fd76d583060cb641f1ae1a",
		"0x3017fb557f04f1c17a6f08631bfa4b08b4cd2236",
		"0x695d34111fa4aa609429d577d09242cbbd14b717",
		"0x4e8985f0b3db04d7d4cbea9e4d2e11233dd65e1d",
		"0xcc15190c103b0d09e69061aa14d9773cb3f744bc",
		"0x48e5b697790fe5746a322899e72fe3f834dffccb",
		"0x5ef4679e46ccbdf9ac9905ffe782dc6227d541ba",
		"0x30101e71fc93a99ba88b0108a5ac305623bb7a55",
		"0xe548e943db1689db91aa291247c4f630d6dc23de",
		"0x13d8267feb280cc6558b87b1570703192e041147",
		"0x799041c167077c9cb60634147a79bbc0372cc94f",
		"0xdf83cd4d2abf71e159fdb0272f38e0b66c4a4a1d",
		"0xab972486abe962cff6705169fdfa8715b1cd18bb",
		"0x1527fa3a0778d574adfe714b7e60c7b2a1bf7c29",
		"0xc9db3be77694f843ca7193aedeb344e3b00b96e8",
		"0x5b25c545fef88be4b1d924330a4393c784519805",
		"0xed33b767cba328e056b6d1a140b7730cd44d7cba",
		"0x5d6b0d1666e6da01631c8d23c7231942d345180f",
		"0x3800d0c11b07adcc63e983e6902ac1c83970c2e5",
		"0x8bcf8878a9b9652f1dd57880f81668bdd5e4f6c4",
		"0x0fd46c0c561ca514ecd6614a08c0063b6fa0d4ca",
		"0xe76c57738cefd5d17291b65a35d21c0d20feea50",
		"0xf5dcb0fbe27629e3f9b7b4991c02ad88edf1e7e1",
		"0x439f4a664c4d5e7a3b3462140d4b820ae1c32fe4",
		"0xd86ab2c9d566b93760ad48e90eccf5f34980e265",
		"0x644628d696fd07fcd8bda2e8b9a680734cc80a93",
		"0xce68328bc32b912aadf7ee3f72d60efb15b7a7dc",
		"0x6128d2eb2daabe8a097451d145d0997028369a47",
		"0x40b8348065e162da1a80d62cdcd1b8e688f8974f",
		"0x772d1fc1b23163f72c08fedbef59f1126cd089c7",
		"0x49a27eabdf13daccf5739fe4741408a2483256e5",
		"0x46436f1fa2f9d77b146cde18d545ae3ed28057c1",
		"0xec2715d333807614d6ae56c1382210b929cbf20f",
		"0x87fefd4dedc92d9ebe9b6dc6699653b7a5517a6c",
		"0xd34427950f74a15bd18d1be4cb3dd5dd65f9eb1f",
		"0x175da058bd3160bc588bb8e442607a8d972286a4",
		"0x823635f24449bac76767afe2afb2616bea3cdf28",
		"0xd7197c91400a7bca1e792536a944d6dea73b8cdb",
		"0x188b2d6ac1a8752ea4ef4adb582646f074096aee",
		"0xd4a36412cde452118b36c8692c5266f8a54d84ed",
		"0x93d7139ff21b2e1b3bd605712c9823d1d9b900e4",
		"0xb5870339663d8a79db881ee168f77d7835b2e721",
		"0x4bba49ec8bbf04282a5bc852176bd79300c06076",
		"0x9e8a7dacce8d640c8748c9f7ba3c647927b6849e",
		"0x645580ccc25dbe305d9c981254cac2c815af4226",
		"0xbddbd3d85d3d14a2d39b5f8121ba7f898cc06c50",
		"0xd4cf1b1793695e115eef3834a32ca175216758c6",
		"0x4a8b9e2c2940fdd39aceb384654dc59acb58c337",
		"0x97b1312006e219e107b8098b9f886cffbc265381",
		"0xe32c56a11b51ed181d6bd87ea6d2233a0a0ac9a0",
		"0x9477d724b5c2e8bf27b4a62795295da76855b4fc",
		"0xfa82a4a0f0cdf4c0f02c938d8840983385c9082d",
		"0x0acc9a662f3204eaa90bc40602a2eed7f909b63c",
		"0x320644dc880390c787d9478d7d529632f33ff9c7",
		"0xf3d72abd69b7c157d3eccc46eeb3ab72b7c85a76",
		"0xe21b790956d8fc776cc8841bdafaeeff414fe717",
		"0x74f0ef88878bad5c38ee76144ebb723c94e25a86",
		"0x9829dd2ea8b55c99875298525c2b5c60c0ca6289",
		"0x26e942673ba611ee6e5e7872ba3020951321e0b8",
		"0x19ea9301b7a47bd0a329723398d75013e7e9ced7",
		"0x4e5c20ed81b3ed1a35e6920bd0ba1581ca174052",
		"0x5aec550cf7e3e883e311ea6d020003dd0afdf6bb",
		"0x746a4b63b86df614bd359b57b1b865531e7bffb0",
		"0xcd3648bc82d6d0e615c92887b18c50ff94414fa0",
		"0x53c61cfb8128ad59244e8c1d26109252ace23d14",
		"0xf4e4587ac84aeb298941d99f8f258bd699e97cbd",
		"0xd33fb5f87481691faed56461de1349bdad55a45e",
		"0x88bf50c76456414b7be74355639769add7b67965",
		"0x36497f6f6de8a4bea1e7bad02c6fe905fb7fb80e",
		"0xf318cbe8aada5e38c6f4f4c7265f7945372d8b83",
		"0xcca71809e8870afeb72c4720d0fe50d5c3230e05",
		"0xdd2bb56f01b5826c99ccebf1925d564ec0d900d7",
		"0x62e41b1185023bcc14a465d350e1dde341557925",
		"0xdda30fef877c4a80e6ed4779a589e34169a87eb6",
		"0x77a0b47b0f7980ac4e80f68c62294fa5177bbfc5",
		"0xb65b3d5d2ebf7eb3e95b88d07617858d0c00f0c3",
		"0x12494ddab76c0004dcc43c0983cc5e48db5eea8f",
		"0x7f8919fde017c5ff8408e4def990bc81f20d3ca4",
		"0xb87f1b304fc557f44316eb28f7563114bb9347f2",
		"0x1e8f169ae605d399e652884e1cc1f6db6fe9550f",
		"0x0f0ed36c79ddb8d774946979864c8b16ddaf2d9c",
		"0xc2c62b99325dc3b6750b2965914abf3fb11b5900",
		"0x57757e3d981446d585af0d9ae4d7df6d64647806",
		"0xbaaad262a7570ab1b5333b9fbc6fd74cb761a91e",
		"0x5167b40679b9fe1657d990334a0ccb4b519514e6",
		"0x8327329a453571ccb74f1c269a4a0acb122e6550",
		"0x5071d7104ee2e1a1c80e620e4198d8bc42b75e17",
		"0xcc031c6c75b42810f1f397e2245003719ec6c031",
		"0x8fa335b8d8db8647cb9e4880ec75597af2c89d8e",
		"0x3859c3404169417bb9b031306cd4a36abc6f4e5c",
		"0x9a87a59040f86b5879cd21cf0b0c42a0c3e2a2ad",
		"0x790a387412d410cf7a184ea8cd162b116aaf0194",
		"0x224fc41805d007775260bcb5f9146ecd5c3d9495",
		"0x09ccd85a14d132551ef0b069a81bd3f766daaea6",
		"0xec6d36a487d85cf562b7b8464ce8dc60637362ac",
		"0x52ffeb82284954274be50f31455648469cca2a1d",
		"0x73da8d83f495a91e3aac0de4d5181f379e1febdf",
		"0x013566174a442a5ed0baa72f77fec2cc8f9268bf",
		"0x1c94c8b5ae2b3efd31629bd252a21540c2be1ff0",
		"0x38cb169b538a9ad32a8b146d534b8087a7fa9033",
		"0xa344b84910b784f1918ae1391971ff9d973902c2",
		"0xaff6bceded11d06b5f554d58d4ab3525fac0284b",
		"0xb541105909a1d80ba28cb4db425026a4731b6cb2",
		"0x63ede6e19401a79889b9cb3ba65f19440084bb11",
		"0xef6e07fe4086a4165149acd1d099856220ff67a8",
		"0xc7b531fdfe0242e1daa47ee2c8efd18fc356b6ae",
		"0x919422d0386726a4dc11b43a65d7501b62dafbe6",
		"0x6d5596d228c44453b0e16d1e4abd3f14724773fb",
		"0x9a7b54749b9a7122a9379d1a244bf45644095907",
		"0xa4cec28d686ee53bc80dd3d79b404cf0dce28016",
		"0x203dcfc573e43735b6d695cd3f43afc57d3fdf19",
		"0xa3b5beb23f46221d7df3a9865174956e89fb4b81",
		"0x3ef72699407af510c0e0351870a6115539c17a70",
		"0x4869e73d9d194f50ab3d76c8107a837b60e6c9b3",
		"0xdf9244945366060ae8ad8907cd08288ecb490b52",
		"0x0ed362471cec2459e82459fbd39c0c736ac38d98",
		"0xc6c47fcbf12ac23ac9efbc496f9baae431fd412b",
		"0x688a2c7ff29d934b3ea0626be78c665ff0ead317",
		"0xee73783d916bb648448ba40d55b9400c781cd374",
		"0x93fe7e05d44b2438c9e2d8c0ce59ccfde8f48422",
		"0xc40dd7acec8ae37b8d7693a5674a567f65a2d486",
		"0xffd0f859442d30938f2bdc788941ad6b44509b3a",
		"0x1ff7e7b6f87fdd13bfa80995117974f9dd9595e6",
		"0x45f2c4e8b4ac24e45e701a13a1c4e89087e9ddbf",
		"0x7846b30044e6ecbbc1b16ecd9e04e78b9e5d3c48",
		"0x549a3ae96da59910fe002519623cb1d0170295b9",
		"0x3548fbc04fc033e07e8004904d18fad502f4f307",
		"0xd16f9281b5d512b480d8b85cea54b4f4dbfa9ff8",
		"0x0ec274b002806af8e19590a52c8a28ac7a40020a",
		"0x8bf7774799f8545bd51a6aa4f306a1b82e50f396",
		"0xb2bb3014af87f04b73304345341fb3d1680d3dc8",
		"0x257f26804f4d6defbfea288f892bb1d687368123",
		"0xd951831d42a7ca028a7da2a6a0da8bb974395665",
		"0x155812dc01f4d24ae29c230bf6505eeaf1547019",
		"0xdd5a0bfd82f6cf9bdec7c9ff23704489c4fd0ccd",
		"0xe91e8e0a72d039df1153297cbf8d8c9e68d31e53",
		"0x1a16c74a408d71fd3314b414114c6eee84aa279e",
		"0x544c3a5364ae3cec04bf8aaf424053aa75e28e85",
		"0x88a8d8a13a69a182004dfa8dc680bee75df4bda5",
		"0xce3066df90ee06331e38d3cbc253d9abafa8e3cd",
		"0x2ccb1cb25695af8f4ff38352c2696a0af8297d01",
		"0xd2921ceebe24a692cd7fc3db6327af823c84345a",
		"0x5b30fa31f5f5b69387ad57c6c33a2d2e8b24ed51",
		"0x758abdc961f0db4f3ef5e077cf4e5c71f32ee5e0",
		"0xa23a02380d55b96fbc6bbefe1c8617fac53a1004",
		"0x4c3cadb49a303bc743308e7a55356ce9f19ab5a4",
		"0x50aa878fb973f91650e94a679c27edb024cf7518",
		"0xedeacf218029a5c434241cfec4c741c8284b34b6",
		"0xaf40fdc468c30116bd3307bcbf4a451a7ebf1deb",
		"0xa2349c899208a82664df9464ae92e46735aba979",
		"0xd1fa0d4cf6c2967578756058b89d076cee963550",
		"0x2569f7c7b250beb3e454dd3b5a00226adc3df932",
		"0x7e593b86ad99fbd60b7702df7bc64d57262c8342",
		"0x0047ee4eca833878bf5c88ec216b1f8b4680b2b3",
		"0x1749ad951fb612b42dc105944da86c362a783487",
		"0x82a07f3d1573cab99da48fb065feb959481fcd22",
		"0x71b672244529e983237c34328ee03134499c222b",
		"0xc445ad9741b8d4626cef03ce8ec3bb3f788085ef",
		"0x78e13f5304581d66421e13208481e5661cc798de",
		"0x171766c36e50f5cd7575b8f40d5ba84967798559",
		"0x5be996fe0bb0af7e7e28bb033c73171b1f9e72d0",
		"0x20fc2324fcf59a6e6d5fee431536e549ff081cc4",
		"0xfd68082c61f208183cf55b309de214206b8f0082",
		"0x097cbe040f6f79aadf6ebb6614d5498e2944df0e",
		"0x2b52775cc13d3817dd908367ad5aa6d27c9beb59",
		"0x5ec54cc2096ca7bfb22c220c875750179cc87faa",
		"0xe27d93738d761b99529aca8271f42d6c8c2e73a2",
		"0x8604d930853c0af2de7b549e03b5d5bf35204e76",
		"0x07649c2241f98ee3779ded23c14459cfeac36b58",
		"0x198ad6c547d20d70f2f656a4f48e6c7cfb7b4325",
		"0xd8b1bb9bccc38b9293e466de144b06f7202f2eb9",
		"0x8971339ec48877b0372ab9b46e63a424d5c1182c",
		"0x30ed824a08389ef85674712bd0203c7ade7bd268",
		"0x5e0c64242fb570b3abd08b7178e0e3d079f5d85f",
		"0x0cfd800dda45deab1e7a03f5b8b4c9e0c1b659bb",
		"0x81a9aed213c74798deec32080b8837ca421614a3",
		"0x3fe0c6d8adce33aeeef51b54e505ee4c4dd39304",
		"0x2dca7bf37c703e685b8ffc2985f6ee6c96a3f9ed",
		"0x1ffd03cac9c17e07addf7eb3255d529bab27789d",
		"0xa86f67136c7d23c0fb21d87dfe4497a1f5bc0c7f",
		"0x3826087e986b68d6142c5258e1f0ced0b6a4d1a9",
		"0x71a0b5d9ecdaa17e15bd364a1b64220ed1d81cf7",
		"0x81016b5fa82b628e7653e63f43882009f90dc2b6",
		"0x8bf4eaacd76ec06a2cac5453071ac4ad3ea5f396",
		"0x7d996e9253c9f491e9993985e9502005cb19b780",
		"0x7673029886e773b91a857f589699dea8dd291ccd",
		"0x591a006e642fa4a18d08702acd5cfe2f3e069e21",
		"0x5acfd581f7c7a243506388a072b3c85e0958e17e",
		"0xa0e585bda18be679e078de5ff50baf627dc64a32",
		"0x5cba1b982ed1c9c5a89b7aaf83556d753b90f2d1",
		"0x5ae98d17b83a041c1b9732423c3f5743b5c8ec0b",
		"0xf1c1a4d8e9312a2e0567cdb276622255510042b7",
		"0x9c06d0b4ae9524e3b726946b04c754b5283fcfd2",
		"0xbca8892f463f35ef20c613e128976f1df114932d",
		"0x2c13e5161abdeae96a56b0112e77c0c45bc1200e",
		"0x3c44151439965c709f7d79ceebaeda5bc5fba9ca",
		"0x934a8a9e6b19a1bd46fa49366d120cb2567dbe6a",
		"0x354f08943fb06a04fd05773d71f34ae1b912bee8",
		"0x5f61bc6690672de640794a8a19f22752ac35632c",
		"0xcc7244296cfef1728118ef36dc9ecc2ecc784a73",
		"0x09bacfb1d3642893ef1faeae2efcda3afb7f6917",
		"0xad95f8c78286c0efc1fb618563d244144a735e12",
		"0x7fce563e9f39ab7dfbd61978f10037cc32a95a4d",
		"0x184813ba43ff87e0dc3bf9fa0d904a4f29457627",
		"0x3fb9d44bc83d0da2902e6230aed42fc219f8a426",
		"0xa5ee51881d9f88f8fd5c87840bbc157c32106efb",
		"0xc0ff12c44c3c581d739ee5eda59567d55d1f01e2",
		"0xcf32352be14b4c39ea7b914f8c302de4d67f0ff6",
		"0x69e0e2b3d523d3b247d798a49c3fa022a46dd6bd",
		"0x50656e03549d2db3944b433a4e8967a7b69cf4a0",
		"0x40db9f106b1cfa1d070f1b1476b1f073f7317186",
		"0x128602b2307798192c92db5c031cfd686a98ee07",
		"0x86a511c19d8942dd8d34ab26ee023502befbe6b3",
		"0xf1ae67235aa3b11851587b9828f9742042ec7713",
		"0x482ee477952006a172ffa7d65346390e046fb83b",
		"0x790c81f69ea3036cac2da74b2d22d6db1e636f6f",
		"0xbc4485554825b3ef94700e179d9e8a98540887a1",
		"0xc0f8ad6e1a7ff71382676b09319d91bbb6344fde",
		"0x449867a1ade267aa9a66e1ee708e3fba578b9524",
		"0x932600b8db3ef4fb0405781ae50194a3ee8fc42a",
		"0x2e62174c79b893b4cfd12b9a3f9c1cca463297c5",
		"0xcbece8ff837ae7b3dc7d1b166779e2cf4050e908",
		"0x98dd5945dc8b039e23658f969e5e1c88b722c771",
		"0xa4b6153679c8ae2e93e3364645b14d3c48c1b142",
		"0xa8d0e92aad42ac33f65f06c3d4f50c71081362a5",
		"0xe9e284277648fcdb09b8efc1832c73c09b5ecf59",
		"0xd35879dcb8b040fd8d6ecc5f808f0b4d3c1cbb0f",
		"0xdb7429084f15dac30e972e5f1f395ef2c268ef3a",
		"0x1b8ef3453949d6b04a6579862fda775870f70ad0",
		"0xb67dde49207ac13faf5fedf81c1671c72241d8d7",
		"0xe664b56636965a3b310e5346007890244ff0c44b",
		"0x7c43c2a60f29186bdaebf320f38350c830ad7fa4",
		"0xa5f4fde5aa9febeb566ad57d08f1169ae54ad5e0",
		"0x9acdb9698f65f2eadcc2813706f83c70a482b3e5",
		"0xe21dbac6bf2f139bc6e16e32541baf308956fbf2",
		"0x4229f57a94492edc977a78ea3016c14df6a199dd",
		"0xf203af2ad57138b6cf24bed16cc73fb80e40f5ca",
		"0xc9b167bc4651518d27d2fad66da487e9a69a6e31",
		"0x09933f32f66de14f2df72d76b307d92be84a8a55",
		"0x3e01bbe365f1d53cfaeadce29102e766d697b4c9",
		"0x98b0059588e47a6cf16c18ee1008d29697badaaa",
		"0x3f7d1580862aa750c26980c944a03ecf57019e9f",
		"0x26e18b9f3a04b18a74577feb9657ba1c925aaea3",
		"0x0167f84d05f56012e8162f0055bb0a184ada8e50",
		"0xebfb007ac8e6239baee43d4812e0373783e0d34f",
		"0xb8d34bf5f55fd764ab7b70c70a3638bd9897c143",
		"0xba1d0958c6589ed64c22144bafb86ef3bae69a60",
		"0x8426ff902eabf20ae698b9e01e77c26885bb3d9a",
		"0xde0f3f1972d5150286220581005e188f66fb3942",
		"0xfff817aea9c28364b3b233e813125ffba1d22e47",
		"0x8643cdfd3197f2c8e7f850920fa6e5b80d958167",
		"0x56568859d2524f13f15a7a12b0a9a76b192c7044",
		"0x162f1b5669e81a5688388badbcda9d632e8f4d98",
		"0x398b7c52ca4fab73dcb51c0e27dba812e29eaede",
		"0x995c8fabcd1259097b8c8a3492ac2b6bcd89dabe",
		"0xa7aad9938043648218aa0512383e32d53c82d878",
		"0x3077583a8c349ed1a84c7d0dab860515a4e25fd0",
		"0xc0c3ffd02ec02128101f0c16f719f623c7d17e34",
		"0xd6893242c511f6773aade60d61e6b07e4f985b6b",
		"0x8d50dcc5916ff342f6752a21cd16de9a4ee3a388",
		"0x9950fe8bb33d71bd020232cf16412f0b17155237",
		"0x6ee5aacca83d658235aab06b7446847bdabf8cea",
		"0xc314eef6afafe3478c4075630c8b613860c2cfe6",
		"0x19c29401044097d450d25038a45898a57a7cea6e",
		"0xeb9c7eea0285e7df0fe34bf83aff6413a474ac33",
		"0x3a25941220ce4ff0c8a41cd9a0d4f66e5548305d",
		"0xd03d26b36642c8137c77ae8fe91e205252db1095",
		"0xdde519c3092c074cfcb2358b7f47c36df9a779f3",
		"0xaf90b5d164dd44160b0b401b2e5161a449d45854",
		"0x284e540c527498d5fa5c48bc852155b0e3ba042d",
		"0xb5ed9734eea4cf42b4028eabee59cd5c82434d05",
		"0x4d8abbeb7acbd33446c8dea413d63e40de375989",
		"0x6537394d5a4919c12d875474a03ff34334c964d4",
		"0xd8799c788ea59933af1c54c17f345408bac67ca7",
		"0x63389a164dbe1fb93d12388b1f38b49203413459",
		"0x1c6e20357c117b91a3663b9918cd3bd731eee604",
		"0x2ed6ec33f7ae5ca9e7256cba0f061498ba97bbb0",
		"0xcc6c1d21e8474b3578e69eb036c712ab08ffdfbb",
		"0x111a92cac4a4c951e3b847909541d07ea325e5b2",
		"0xe539db836291b6208489f74f54d5935e088f8750",
		"0x4e88b72d9d312b729c1d9e164d976de14d00bcd8",
		"0xc2ae44253c5328189b4c4a5664d14b474ad3e08f",
		"0x80c4fc6fd31d6d65dac5a1781f8d177a5ded69f6",
		"0xb5e61ecdd89f5cba7e6f9657329413a40147f958",
		"0xfc69130bd08bd266894382eef5ca18631e1f58db",
		"0x5abd2d5e2430af291f334824c93181583fb82c05",
		"0x2380e45eb243f685424c056543641fde4f8c7635",
		"0xfa5233de3db673d48a1bb742c77180c0b83049d3",
		"0x52ba527bf60bca4c5d49b034465bb6e7c3494206",
		"0xcdba19e0e34d6ac31922b1c153d1f6c80bb03a60",
		"0x06de102ff2395820a65340ecf19c9da370310053",
		"0x895607a6f97678d15e72c73728518b1fd0729c4b",
		"0xa33ef019ff0bde98f951585aa57c1af7c898fe5c",
		"0x9a96e82599ff52831b03e2c7245bd8e83b867618",
		"0xc925f8636f37d8060a848804f96fbc6b4d527d06",
		"0xf7b2ee9174c0c347b73bb61d98b83f708c35d3a1",
		"0x9674a1578d3519bb0051e0cbfddaed642c12309b",
		"0x815dc6d97506e94384fa350678cbe2c2c30c0b40",
		"0xcd5b1e9a0db0a31f7c46fe2e4a3d204e61f6dfc1",
		"0xa80e26a97d7758abcc7a8bbcb42ac8552e0919c5",
		"0xeba40f475728ef05469f90c831f21643fe3656f9",
		"0xd71f271ffa7a6580e5710feb7279d795a8272a66",
		"0x8e812e03cc6b2cefefac00a11551a17153592743",
		"0xa185b29e007e61f87b67e60b9515132a529c5fcd",
		"0x630a913a9031c9492abd4c41dbb15054cfec4416",
		"0x84bf0c0a3f8f7150fabbc0111f77c5d6e9c1bf8c",
		"0xf685b2fc0a1372043d3f1dfc73f721613c8933b2",
		"0x298afe764ceade43ba80ac4390195e8d9451d6ef",
		"0xd64dd17b4aca744ac5830cc4a2f5e7165300e8a1",
		"0xbf0dbb594eabbe072c5370d37939a5268fc676e9",
		"0x6efec80290cf686c79d4df05f7d60a4412913d4a",
		"0xac4d9c0a82108ed85129c8f86a49aea8edeeeccf",
		"0xdcd62b41aabf29a2249180a5efc7ab8e974de23b",
		"0x446bf052bf571769e768e172da02b7e6af0bfa92",
		"0xe1e534f25f2cb54f8d7de306fd936bbd6091a994",
		"0xbe591b1bd3e5d4ebc04c71552611dde6864b26a2",
		"0xc224ebfbb4e3bdfca4134b3697313984a8936b5b",
		"0x921a1cd96363578f52c0d7cc3034056b46e0b30e",
		"0x166da02e436855dd154febe2e77d9fffa193156b",
		"0x764ef1cda89af77c4458ea1b70fc65c48f3f7954",
		"0x98f8407c1dd18d6ea2ecc48fdf772665ac5323c1",
		"0x2b988d1a98ce8f583ae81f234f8af83c0b50fd79",
		"0x262757e7f40c592e990b85623b73c65afa6e53a8",
		"0xb6755293c29452cc9a709a303e371f10e68eddfa",
		"0x4ddd3c40835736df066ca70d0e8a0346b1366673",
		"0x147e8a4d0b2f5b3567ac3777ec640def30c17d90",
		"0x805ba71c611f791b428e527b6d6cc801589ae1a9",
		"0xa469418efabede2452c813bfc85c3e72e987dad7",
		"0xd16833e93a17ba0a241d68d08c125c8acbe585f8",
		"0xa289364347bfc1912ab672425abe593ec01ca56e",
		"0x652129be3f15650b0cd22b74706cd6cbcdffad41",
		"0x61146f999ffdf7b260b239dcd1ad1dec4733529b",
		"0x5c31d69f2312c37bedec59624481ea3340edcb7b",
		"0x90dcf653aa3302a9caa98a3bd632159f257a8520",
		"0xa9453563c85565ae94cc784340f95fde3151244b",
		"0xb4a40a766547a628594dd25d2e794ed7cf03f0d1",
		"0x7feb4b824e3c1efae8223b721f2588dda701a8ee",
		"0x1c7920062877e987a0b926d0d2b92013b5b3326f",
		"0xa9b3a91e058f820ed6f385a6026478de906664ea",
		"0xb3ad178d923a26c74e99d8bf2a5b6cf900376d85",
		"0x89febd02c6f54c766b586de8335e04958ae3c13a",
		"0xc50c40903908493e7975c83c809e1d3853a10e62",
		"0xff5f090c38e698f740a109148793143b2c8caae2",
		"0x0c5fd74e11ecf0bc2cc0739c721714b5e25928ad",
		"0x78de1415323168fcfa375c308b8017e2ea952615",
		"0x49a4f0bcb3e20d59ea0456e19eec0a15134a7e50",
		"0xbf01e71bcfa6f01ff0e1a2ca73c1109bd15625a4",
		"0x4d17676309cb16fa991e6ae43181d08203b781f8",
		"0xed05860e99c8063138c531d3f058fcd5d777036c",
		"0x3413f5719464c66b434532c3b844301dfb437daf",
		"0x4075c36ab1895149435031ff8ef9d3e52ad01ff8",
		"0x0148107d2b0fc3e90e8ac427a9058beaffc5dd79",
		"0x9421130eaa5acc1f5f3a03eabd5baf37331f2c31",
		"0x4e23e93d9fb1982d194ce64b4dae629c086539d0",
		"0x0dc966e43f5dd38e626c3e0cd5cfb7c17bd8e0ff",
		"0x52be13f3fb879f2004a3caa14562dcfccf8d0441",
		"0x0df6971b32cb47198946a1298431dac63f25a9bc",
		"0x5535614abcda341b0d5708f7d34bbe822331f841",
		"0xcf077979998923ca3ed6953e58a41b1f3c667cd2",
		"0xe9b9a2747510e310241d2ece98f56b3301d757e0",
		"0x5d6152ec2c222160a41749ad9c19649df8e27594",
		"0x954ad434e1490b97a37966709f3de0e50a252c91",
		"0xb16a0eb9bfd54311e0c7073f99ac7503f19bedee",
		"0x02b322ecfd80aa1de6536be93a496820a84711ed",
		"0xd571e4efa97a186889d65d921c726781cc6ef04a",
		"0xfd4f24676ed4588928213f37b126b53c07186f45",
		"0x019dcc093a47fb23ab7224d9b4fdbe1944d012c0",
		"0xfb626333099a91ab677bcd5e9c71bc4dbe0238a8",
		"0x2a5d13655ae7cc2150793d795ea8122ac5c46d7b",
		"0xf6abd8fec49ea77131997a22d3816427b0038132",
		"0x8307ca61984dce4dc5cd8645434a5ad2a82cfff4",
		"0x1244d96845d94b3d1a550110a5c97ef018f4d24e",
		"0xc659ba9e05a6c72a946328d528b3b8bcd4ebe98f",
		"0x042345ac4de6d983750563cebd2bd8ce4c802624",
		"0x3e7eb50bc57c63ad39ede448d7851d96b5b2a22b",
		"0xac8433e5337a0dd7346c1678bea4d1f02a068bbd",
		"0xf27a63071354a265a8c7ee89c3ca9cb4660ac170",
		"0x8626354048f90faafc212c07ec7f3613406b1b32",
		"0x77e2a981df51d10c964343fb0607cf6112f034ef",
		"0xa021ca2f5ec4dc539655718bca4fca5d0015163b",
		"0xbb608fc161304b5fbb7964aba309d42eef8a5835",
		"0x5a537f211709f95bab041342bf0067cddde7ba14",
		"0x9b4eb15f18e5d1bc4211b8b8d521c534f15e7a3d",
		"0xb609cf81ef84bae686c719764811feb34e12f7f8",
		"0x602f2e120a9956f2ad1ce47ced286fcefbba9f8c",
		"0x0a6ef7ee49df7ebce0cc76a1f7f53ab261f2fc83",
		"0xdca1bd9a3fc69a18e268c26a9691245cf36438cf",
		"0xcee44afa1b9c7387386f203afdfb665c6d2ff198",
		"0xbf9c66d304fe3bd34a54ad9634c999b1f7163522",
		"0xd4eb9067111f81b9baabe06e2b8ebbadaded5daf",
		"0xc0360b02c1871c27a8815f6bc094ec58557c8d50",
		"0xe87bb1d31e8b69a51ef38b5e73dff2efc32fd309",
		"0x6f61832a0616bd8625a3e894d1d65a51bb17f503",
		"0xa1d3c765e9a9655e8838bc4a9b16d5e6af024321",
		"0x04535af928fa83fe7a8a1c3edbd464d0f3bc14ea",
		"0xa09450e986236bd6aa5f211c07ef7a7ffeb4a0a9",
		"0x12e8b4e8964164a556dc1b32c3af2aa5a61a7af8",
		"0x3e5fd0244e13d82fc230f3fc610bcd76b3c8217c",
		"0xfd9deefc3442d5d2c9ebe4470b026372ebc2a948",
		"0xd98d31b80c29df979c5f10f89f50f9a530cecb96",
		"0x38abd5912b5e2537b9e68abbd8d037d0e35eba0e",
		"0xa238d8b887d7a9cd11d26b59879480baaff33c84",
		"0x3ab55f46cb8b5e4966849a3f189053fdb7190003",
		"0x0e0a1078a7dc8a4f564b4ebceb2c6e158980f2bb",
		"0x78992d75022100d2d0c5c7fbf3181a1f92000813",
		"0x6169d817bc267a8adfb4229a39e06a125df28a5c",
		"0x2e9f665ff67a93ed66780b5265fe24e2db8b3c85",
		"0x9b4bd6b3f1fcbde47c7222e2335ff12a9a6dc1c9",
		"0x54f252f92d58bd5c3a31c40b71cab59b8bc85cd4",
		"0xa66749c5012ca1377b92a563474b5c598f876b66",
		"0x77c9fc9729b848be2bbc3fb1a10a31a7e1a31765",
		"0x51c260aa98b876a49733375696cecb15df0a9d7c",
		"0x958681ea23a5b9f77fb5fbe3e054e5a8831f412b",
		"0x713dcf1ddc4ffc859f78a7dcfbe2d2bebc724831",
		"0xec8abfd5cce0198914d76c733b5a0763a88da24f",
		"0xa81ace214b97d4b9c2072a934d0c4de486757538",
		"0x6aa35f8ae7c9aceb2a541b6576577fdcdae73100",
		"0xf8a0ceb25548aeefd558d1d9b5f839380b1f3485",
		"0x7fbbafe8f03809d3fd957a56c9c882d1566f88a1",
		"0x822ca1fa448f8ff4018650570c6bc709e3dff846",
		"0xaf23f3936b14046b7388fd8b6af4522525672b72",
		"0x87bcfdcc9899afd70abe27f4838b9d43356affff",
		"0x3b92707f9e1a958f5c597653aaca43ef06c04854",
		"0x29906f250d7a0f419ab31ae2c1a0ca9663924b64",
		"0x7e2930f87831987a38b302e1169de500409092f8",
		"0x0bd851645d8c919228e142a404bc0b33a78b85d6",
		"0x4586e8eb86cb856b090881600754b70c75525588",
		"0x983f865d02b84134eaeb50b97d7cc5021fb5e2e5",
		"0xeac54ff283b7513e6e369c1432b187ee40e42c8f",
		"0x78adc70fcc9c0cf5f8c996f02fde3d9ddb669d57",
		"0x2b47fe41b81c52537c5555187a29e734dcf809a6",
		"0xb8590a1b9e3f2571c494a66238b3a3b8867f12ed",
		"0xead8dab4e0beeea9e47b69b1c7b89c2db00be162",
		"0xd70300d3d31a38f2c1b302e0e90136a4d24ad972",
		"0xbe9e65a494b8c378d05cc4dcd6b8a996219f7d51",
		"0x2ef2d6e5d34b5f3b465adfa1b38204dbd274b32a",
		"0x7d53f718c8799110bd8598f9aad4a1fb109aa0ab",
		"0xb2aa34a00047012e759dc162c3ce8130b6a1961c",
		"0xa088b8291f7983a6e7915018e94b6a4a25ca668f",
		"0xce42353fd66097a9a87885fa1e66c70c20a0ee60",
		"0x09ae5a2bc9571b7f45a6c1bfdeca9abd2ab623e3",
		"0x7984b13b109a8dfc8a0fdacdfef17e131dcb57ba",
		"0x2944a48428f451564ec40c1e7f354bd9abc6dfa8",
		"0x4a8b3e09638f3aa054b5067cb77447d98c4e968f",
		"0x2d00dea942f7eab75cfd97de5e2395ba3332c9aa",
		"0x726d84e8c399549e40aea16dbae23b72dd89f5a6",
		"0xb519d154698f8a6e5396325ed3b3d220bb513296",
		"0x5839dd13ad2c78ddcf3365d54302d65764619737",
		"0xf8bd0881bdad5470a1d17239a61675fbc98ef7d7",
	},
	"stakefish": {
		"0x61c808d82a3ac53231750dadc13c777b59310bd9",
		"0x4c7052546a9e38a72b9731e334d8100b45e88cbc",
		"0x1fa953ef0c089b342656194e569d2bc0de698add",
		"0x4e017a192177abdae99250b6a0809dd76c856767",
		"0xaa5104d696a3c8ab61331c87d9294c0ae1f5aa51",
		"0x0194512e77d798E4871973d9cB9D7DDFC0fFd801",
	},
	"stkr": {
		"0x4069d8a3de3a72eca86ca5e0a4b94619085e7362",
	},
	"bitfinex": {
		"0x2b1df729083f6416861445d8aaac04ebdcd4a848",
	},
	"okex": {
		"0x5a0036bcab4501e70f086c634e2958a8beae3a11",
	},
	"stakewise": {
		"0x102f792028a56f13d6d99ed4ec8a6125de98582a",
		"0x5fc60576b92c5ce5c341c43e3b2866eb9e0cddd1",
		"0xAD4Eb63b9a2F1A4D241c92e2bBa78eEFc56ab990",
	},
	"piedao": {
		"0x66827bcd635f2bb1779d68c46aeb16541bca6ba8",
	},
	"wexexchange": {
		"0xbb84d966c09264ce9a2104a4a20bb378369986db",
	},
	"kucoin": {
		"0xd6216fc19db775df9774a6e33526131da7d19a2c",
	},
	"poloniex": {
		"0x0038598ecb3b308ebc6c6e2c635bacaa3c5298a3",
	},
	// See rocketpool.go
	"rocketpool": {
		"",
	},
	"cream": {
		"0xb2c43455ee556dea95c0599b0d3f7d0abbf32fdb",
		"0x197939c1ca20c2b506d6811d8b6cdb3394471074",
	},
	"defi1": {
		"0x234ee9e35f8e9749a002fc42970d570db716453b",
	},
	"defi2": {
		"0x711cd20bf6b436ced327a8c65a14491aa04c2ca1",
	},
	"defi3": {
		"0xfb626333099a91ab677bcd5e9c71bc4dbe0238a8",
	},
	"whale1": {
		"0x3b436fb33b79a3a754b0242a48a3b3aec1e35ad2",
	},
	"whale2": {
		"0xa8582b5a0f615bc21d7780618557042be60b32ed",
	},
	"whale3": {
		"0xa3ae668b6239fa3eb1dc26daabb03f244d0259f0",
	},
	"whale4": {
		"0x152c77c1131eb171b5547f2c4c33ff2c496984f1",
	},
	"whale5": {
		"0x38bc2db2732a7efa43c89160118268ef0d15163d",
	},
	"whale6": {
		"0x0ec24e397a3592da41264209cc0582bc3bbc9776",
	},
	"whale7": {
		"0x46d24de64abe6a7ab6afda0389074ae8c7ed5d0f",
	},
	"whale8": {
		"0x6aa35f8ae7c9aceb2a541b6576577fdcdae73100",
	},
	"whale9": {
		"0x27d71a464da941118a4e056d0737274aa308a923",
	},
	"whale10": {
		"0x51226113a35a29cda0030ef221e90baef221278b",
	},
	"lighthouse-team": {
		"0xFCD50905214325355A57aE9df084C5dd40D5D478",
	},
	"prysm-team": {
		"0x7badde47f41ceb2c0889091c8fc69e4d9059fb19",
	},
	"teku-team": {
		"0x43a0927a6361258e6cbaed415df275a412c543b5",
	},
	"nimbus-team": {
		"0x5efaefd5f8a42723bb095194c9202bf2b83d8eb6",
	},
}
