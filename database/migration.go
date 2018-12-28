package database

import (
	"time"
	"udabayar-go-api-di/model"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Database struct {
	*gorm.DB
}

func (db Database) Migrate() {
	db.AutoMigrate(
		&model.User{},
		&model.Profile{},
		&model.Shop{},
		&model.Address{},
		&model.Transaction{},
		&model.Order{},
		&model.Product{},
		&model.ProductVariant{},
		&model.ProductImage{},
		&model.ProductReview{},
		&model.Category{},
		&model.Storefront{},
		&model.Bank{},
		&model.UserBank{},
		&model.Courier{},
		&model.ProductDetail{},
		&model.ShopCourier{},
		&model.PaymentMethod{},
		&model.ProductStorefront{},
		&model.Testimonial{},
	)
}

func (db *Database) Seeder() {
	// User (dummy data)
	user := model.User{
		Email:    "dickyindra@gmail.com",
		Password: "$2a$14$nMRtfDlp7tLTGNEgjl1LgOWFCIHOAX532wvCMc3c0YELfrZpsxwzq",
		Name:     "Dicky Indra Asmara",
		Status:   1,
		Profile: model.Profile{
			Phone: "6285292365555",
			Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTj-uzauQNtOhsvk3Aegrl08ynRe3F5wpt9WGtq6xl2-ANkKMkP",
		},
	}

	db.Create(&user)

	// Testimonial (dummy data)
	testimonial := model.Testimonial{
		UserID:      1,
		ShopID:      1,
		Testimonial: "Pelayanannya bagus, memudahkan para pembeli untuk melakukan pembayaran.",
	}

	db.Create(&testimonial)

	// Bank
	banks := []model.Bank{
		{
			Name: "Bank Central Asia",
			Code: "BCA",
			Logo: "https://s3.amazonaws.com/udabayarmedia-bucket/images/banks/bca-logo.png",
			PaymentMethod: []model.PaymentMethod{
				{
					Name:    "ATM BCA",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan kartu ATM dan PIN BCA anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih menu Transfer &gt; Ke Rek BCA.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nomor rekening Toko (sesuai dengan yang muncul pada website) dan pilih Benar.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nominal total tagihan dan pilih Benar.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Periksa data di layar, pastikan nama, nomor rekening Toko dan total tagihan sudah sesuai. jika sudah benar, pilih Ya.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Bukti pembayaran akan tercetak setelah transaksi berhasil.</span></li> </ol>`,
				},
				{
					Name:    "Internet Banking KlikBCA",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan user ID dan PIN internet banking anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih menu Transfer Dana &gt; Daftar Rekening Tujuan &gt; Rekening BCA dan klik Kirim. (Jika nomor rekening belum terdaftar, jika sudah lanjut ke langkah nomor 4)</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nomor rekening Toko (sesuai dengan yang pada website), masukkan angka yang ada pada aplikasi internet banking ke KeyBCA, dan masukkan respon KeyBCA anda dan klik Kirim.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Kemudian pilih menu Transfer Dana &gt; Transfer Ke Rek. BCA &gt; Dari Daftar Transfer.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih nomor rekening Toko dan masukkan nominal total tagihan.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan angka yang ada pada aplikasi internet banking ke KeyBCA dan masukkan respon KeyBCA anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih transfer sekarang dan klik Lanjutkan.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Periksa data di layar, pastikan nama, nomor rekening Toko dan total tagihan sudah sesuai. jika sudah benar, klik Kirim.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Halaman konfirmasi transaksi berhasil akan muncul dan dapat disimpan sebagai bukti pembayaran.</span></li> <li style="font-weight: 400;"></li> </ol>`,
				},
				{
					Name:    "m-BCA (BCA Mobile)",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan kode akses anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih menu m-Transfer &gt; Antar Rekening.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih nomor rekening Toko (sesuai dengan yang muncul pada website) dan masukkan nominal total tagihan dan klik Kirim.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Periksa data di layar, pastikan nama, nomor rekening Toko dan total tagihan sudah sesuai. jika sudah benar, klik OK.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan Pin m-BCA anda dan klik OK.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Halaman konfirmasi transaksi berhasil akan muncul dan dapat disimpan sebagai bukti pembayaran.</span></li> </ol>`,
				},
			},
		},
		{
			Name: "Bank Rakyat Indonesia",
			Code: "BRI",
			Logo: "https://s3.amazonaws.com/udabayarmedia-bucket/images/banks/bri-logo.png",
			PaymentMethod: []model.PaymentMethod{
				{
					Name:    "ATM BRI",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan kartu ATM dan PIN BRI anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih menu Transfer &gt; BRI.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nomor rekening Toko (sesuai dengan yang muncul pada website) dan pilih Benar.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nominal total tagihan dan pilih Benar.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Periksa data di layar, pastikan nama, nomor rekening Toko dan total tagihan sudah sesuai. jika sudah benar, pilih Ya.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Bukti pembayaran akan tercetak setelah transaksi berhasil.</span></li> </ol>`,
				},
				{
					Name:    "Internet Banking BRI",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan username, password, dan validation anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih menu Transfer &gt; Transfer sesama BRI.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nomor rekening Toko (sesuai dengan yang pada website) dan masukkan nominal total tagihan dan klik Kirim.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Periksa data di layar, pastikan nama, nomor rekening Toko dan total tagihan sudah sesuai. jika sudah benar, klik Kirim.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan password, m-Token anda dan klik Kirim.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Halaman konfirmasi transaksi berhasil akan muncul dan dapat disimpan sebagai bukti pembayaran.</span></li> </ol>`,
				},
				{
					Name:    "Mobile Banking BRI (BRI Mobile)",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih menu mobile banking BRI &gt; Transfer &gt; Sesama BRI.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih jenis transfer, masukkan nomor rekening Toko (sesuai dengan yang muncul pada website) dan masukkan nominal total tagihan dan klik Ok.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan PIN mobile banking anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Apabila pesan konfirmasi untuk transaksi menggunakan SMS muncul, pilih</span> <span style="font-weight: 400;">Ok. </span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Status transaksi akan dikirimkan melalui SMS dan dapat digunakan sebagai bukti pembayaran.</span></li> </ol>`,
				},
			},
		},
		{
			Name: "Bank Negara Indonesia",
			Code: "BNI",
			Logo: "https://s3.amazonaws.com/udabayarmedia-bucket/images/banks/bni-logo.png",
			PaymentMethod: []model.PaymentMethod{
				{
					Name:    "ATM BNI",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan kartu ATM dan PIN BNI anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih menu Lain &gt; Transfer &gt; Dari Rekening Tabungan &gt; Ke Rekening BNI.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nomor rekening Toko (sesuai dengan yang muncul pada website) dan pilih Benar.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nominal total tagihan dan pilih Benar.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Periksa data di layar, pastikan nama, nomor rekening Toko dan total tagihan sudah sesuai. jika sudah benar, pilih Ya.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Bukti pembayaran akan tercetak setelah transaksi berhasil.</span></li> </ol>`,
				},
				{
					Name:    "Internet Banking BNI",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan UserID dan password anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih Transaksi &gt; Info &amp; Administrasi Transfer &gt; Atur Rekening Tujuan &gt; Tambah Rekeningku Tujuan dan klik Ok.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih Kode Network &amp; Bank : Transfer Antar Rek. BNI.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nomor rekening Toko (sesuai dengan yang pada website) dan klik Lanjutkan.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan data lainnnya dan klik Lanjutkan.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Cek detail konfirmasi. Pastikan nama pemilik Toko, nomor rekening Toko dan nominal total tagihan anda sudah benar.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan PIN BNI e-secure anda dan klik Proses.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih Transaksi &gt; Transfer ? Transfer Antar Rek BNI.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih rekening tujuan sebagai rekening yang telah disimpan. Masukkan nominal total tagihan dan klik Lanjutkan.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Cek detail konfirmasi. Pastikan nama rekening Toko, nomor rekening Toko dan nominal total tagihan anda sudah benar.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan PIN BNI e-Secure dan klik Proses.</span></li> </ol>`,
				},
				{
					Name:    "Mobile Banking BNI (BNI Mobile)",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan UserID dan MPIN anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih Antar Rekening BNI &gt; Transfer.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih rekening tujuan &gt; Input rekening baru</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nomor rekening Toko (sesuai dengan yang pada website) dan klik Lanjut.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">masukkan nominal total tagihan dan klik Lanjut.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Periksa detail konfirmasi. Pastikan nomor rekening toko dan nama penerima sudah benar. Jika benar, masukkan password transaksi dan klik Lanjut.</span></li> </ol>`,
				},
			},
		},
		{
			Name: "Bank Mandiri",
			Code: "Mandiri",
			Logo: "https://s3.amazonaws.com/udabayarmedia-bucket/images/banks/mandiri-logo.png",
			PaymentMethod: []model.PaymentMethod{
				{
					Name:    "ATM Mandiri",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan kartu ATM dan PIN BNI anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih menu Lainnya &gt; Transfer &gt; Ke Rekening Mandiri.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nomor rekening Toko (sesuai dengan yang muncul pada website) dan pilih Benar.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nominal total tagihan dan pilih Benar.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Periksa data di layar, pastikan nama, nomor rekening Toko dan total tagihan sudah sesuai. jika sudah benar, pilih Ya.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Bukti pembayaran akan tercetak setelah transaksi berhasil.</span></li> </ol>`,
				},
				{
					Name:    "Internet Banking Mandiri",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan UserID dan password ib mandiri anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih Transfer &gt; Rekening Mandiri.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nominal total tagihan dan nomor rekeningku Toko (sesuai dengan yang ada pada website).</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih tanggal transfer sekarang dan klik Lanjutkan.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Cek detail konfirmasi. Pastikan nama pemilik Toko, nomor rekening Toko dan nominal total tagihan anda sudah benar.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan respon token dan klik Lanjutkan.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Halaman konfirmasi transaksi berhasil akan muncul dan dapat disimpan sebagai bukti pembayaran.</span></li> </ol>`,
				},
				{
					Name:    "Mobile Banking Mandiri (Mandiri Online)",
					Context: `<ol> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan UserID dan PIN anda.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih Transfer &gt; Ke Rekening Mandiri.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Pilih nomor rekening Toko (sesuai dengan yang pada website).</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Masukkan nominal total tagihan, pilih transfer sekarang dan klik Lanjut.</span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">Periksa detail konfirmasi. Pastikan nomor rekening toko dan nama penerima sudah benar. Jika benar, klik Kirim </span></li> <li style="font-weight: 400;"><span style="font-weight: 400;">masukkan MPIN dan </span><span style="font-weight: 400;">Halaman konfirmasi transaksi berhasil akan muncul dan dapat disimpan sebagai bukti pembayaran.</span></li> </ol>`,
				},
			},
		},
	}

	for _, bank := range banks {
		db.Create(&bank)
	}

	// User Bank
	userbanks := []model.UserBank{
		{
			OwnerID:        1,
			BankID:         1,
			AccountName:    "",
			AccountNumber:  "122301000511561",
			Username:       "Sumardi4506",
			Password:       "123456",
			Balance:        0,
			Branch:         "SCBD",
			CityOrDistrict: "Jakarta Selatan",
		},
		{
			OwnerID:        1,
			BankID:         2,
			AccountName:    "",
			AccountNumber:  "122301000511561",
			Username:       "robbys0525",
			Password:       "Bfx168899",
			Balance:        0,
			Branch:         "SCBD",
			CityOrDistrict: "Jakarta Selatan",
		},
	}

	for _, userbank := range userbanks {
		db.Create(&userbank)
	}

	// Address (dummy data)
	address := model.Address{
		Name:           "Ruko Alila",
		Phone:          "6285292365555",
		Address:        "Jalan Cempaka Raya No.3, RT.02 / RW.04, Cempaka Putih Barat, Cempaka Putih, RT.2/RW.4, Cemp. Putih Bar., Cemp. Putih, Kota Jakarta Pusat, Daerah Khusus Ibukota Jakarta 10520",
		Province:       "DKI Jakarta",
		CityOrDistrict: "Jakarta Pusat",
		SubDistrict:    "Cempaka Putih",
		OriginCode:     151,
		PostalCode:     53142,
		OwnerID:        1,
	}

	db.Create(&address)

	// Couriers
	couriers := []model.Courier{
		{
			Name: "Jalur Nugraha Ekakurir",
			Code: "JNE",
		},
		{
			Name: "Titipan Kilat",
			Code: "TIKI",
		},
	}

	for _, courier := range couriers {
		db.Create(&courier)
	}

	// Shop (dummy data)
	shop := model.Shop{
		OwnerID:           1,
		Name:              "Alila Shop",
		SlugName:          "alila-shop",
		Description:       "Alila Shop menawarkan rangkaian eksklusif produk fashion mulai dari topi hingga sepatu.",
		ProfilePicture:    "https://hijabalilaolshop.yukbisnis.com/applications/upload/business/2015-10/hijabalilaolshop/albums/profile/hijabalilaolshop-logo.jpeg",
		BackgroundPicture: "",
		VideoUrl:          "",
		AddressID:         1,
		ShopCouriers: []model.ShopCourier{
			{
				CourierID: 1,
			},
			{
				CourierID: 2,
			},
		},
	}

	db.Create(&shop)

	// Category
	categories := []model.Category{
		{
			Name: "Sepatu Anak-anak",
		},
	}

	for _, category := range categories {
		db.Create(&category)
	}

	// Product (dummy data)
	products := []model.Product{
		{
			ShopID:      1,
			CategoryID:  1,
			Name:        "SEPATU ANAK ADIDAS SLIP ON KIDS 30-35 SEPATU ANAK ADIDAS CLOUDFOAM",
			SlugName:    "sepatu-anak-adidas-slip-on-kids-30-50-sepatu-anak-adidas-cloudfoam",
			Description: "<br>GESER GAMBAR WARNA BARU RESTOCK !!!!<br><br>GAMBAR TERAKHIR BARANG YANG BARU RESTOCK <br><br>1. REAL PICTURE 100%<br>2. SIZE 31-35<br>4. HARGA TERMURAH SEONLINESHOP<br><br><br>PANDUAN SIZE :<br>30 : 19cm<br>31 : 19,5 cm<br>32 : 20 cm<br>33 : 20,5 cm<br>34 : 21cm<br>35 21,5 cm",
			NonPyshical: 0,
			Condition:   0,
			MainSku:     "SAASOK3035",
			Weight:      500,
			ViewCount:   0,
			IsActive:    1,
			ProductDetail: []model.ProductDetail{
				// {
				// 	Qty:   50,
				// 	Price: 50000,
				// },
				{
					ProductVariant: model.ProductVariant{
						Sku:  "SAASOK3035_001",
						Name: "Merah",
					},
					Qty:   50,
					Price: 50000,
				},
				{
					ProductVariant: model.ProductVariant{
						Sku:  "SAASOK3035_002",
						Name: "Putih",
					},
					Qty:   50,
					Price: 75000,
				},
			},
			ProductImage: []model.ProductImage{
				{
					ImageUrl: "https://s3.amazonaws.com/udabayarmedia-bucket/images/products/39b70be715e5548d2471db06a4162c3f.jpg",
				},
				{
					ImageUrl: "https://s3.amazonaws.com/udabayarmedia-bucket/images/products/a61e5762db63389d66b6936def85df20.jpg",
				},
			},
			ProductReview: []model.ProductReview{
				{
					TransactionID: 1,
					Rating:        4,
					Review:        "Barangnya bagus gan, hehe.",
				},
			},
		},
		{
			ShopID:      1,
			CategoryID:  1,
			Name:        "Steam Wallet $5",
			SlugName:    "steam-wallet-5",
			Description: "<p><br>-IDR 700,000 = Rp 690,000<br>-IDR 1400,000 = Rp 1370,000</p><p>Tambahkan deposit ke steam wallet Anda untuk pembelian permainan apapun di Steam, item Dota 2, item Team Fortress 2 atau dalam permainan yang mendukung transaksi Steam.</p><p>Add funds to your Steam account for the purchase of any game on Steam, Dota 2 items, Team Fortress 2 items or within a game that supports Steam transactions.</p><p>SYARAT &amp; KETENTUAN :<br>- Kode akan dikirimkan melalui kolom resi pada saat setelah pembayaran (1x24 Jam)<br>- Mohon kode yang diterima segera digunakan<br>- Mohon jangan memberitahukan kode tersebut kepada siapapun<br>- Kode yang sudah dikirimkan tidak dapat ditukar atau dikembalikan dengan alasan apapun juga<br>- Kami menjamin keabsahan dan kerahasiaan kode tersebut, namun kami tidak bertanggung jawab atas penyalahgunaan kode tersebut&nbsp;</p><p>&nbsp;</p>",
			NonPyshical: 1,
			Condition:   0,
			MainSku:     "SW5AS",
			Weight:      0,
			ViewCount:   55192,
			IsActive:    1,
			ProductDetail: []model.ProductDetail{
				{
					Qty:   1250,
					Price: 73000,
				},
			},
			ProductImage: []model.ProductImage{
				{
					ImageUrl: "https://s3.amazonaws.com/udabayarmedia-bucket/images/products/39b70be715e5548d2471db06a4162c3f.jpg",
				},
				{
					ImageUrl: "https://s3.amazonaws.com/udabayarmedia-bucket/images/products/a61e5762db63389d66b6936def85df20.jpg",
				},
			},
		},
	}

	for _, product := range products {
		db.Create(&product)
	}

	// Transaction (dummy data)
	transaction := model.Transaction{
		Code:                    "7EbxkrHc1l3Ahmyr",
		Status:                  3,
		ShopID:                  1,
		Amount:                  60467,
		UserBankID:              1,
		PaymentExpiredAt:        time.Now().Local(),
		Courier:                 "JNE",
		CourierType:             "YES",
		PaymentCode:             467,
		PriceShipment:           10000,
		EstimatedShipping:       "1-2",
		TrackingNumber:          "12935829173982",
		Note:                    "Packing aman ya gan, terimakasih.",
		CustomerName:            "Franseda Indra",
		CustomerEmail:           "franseda99@gmail.com",
		CustomerPhoneNumber:     "085292365555",
		CustomerAddress:         "Jln. Moch. Yamin Gg. IX Karangpucung",
		CustomerProvince:        "Jawa Tengah",
		CustomerCityOrDistrict:  "Purwokerto",
		CustomerSubdistrict:     "Purwokerto Selatan",
		CustomerDestinationCode: 1145,
		PostalCode:              53142,
		Order: model.Order{
			ProductID:       1,
			ProductDetailID: 1,
			Qty:             1,
			Price:           50000,
		},
	}

	db.Create(&transaction)
}

func (db Database) Down() {
	db.DropTable(
		&model.User{},
		&model.Shop{},
		&model.Address{},
		&model.Transaction{},
		&model.Order{},
		&model.Product{},
		&model.ProductVariant{},
		&model.ProductImage{},
		&model.ProductReview{},
		&model.Category{},
		&model.Storefront{},
		&model.Bank{},
		&model.UserBank{},
		&model.Testimonial{},
	)
}

func DatabaseMigration(database *gorm.DB) *Database {
	return &Database{
		database,
	}
}
