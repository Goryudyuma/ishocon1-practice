package main

import "log"
import "fmt"

// Product Model
type Product struct {
	ID          int
	Name        string
	Description string
	ImagePath   string
	Price       int
	CreatedAt   string
}

// ProductWithComments Model
type ProductWithComments struct {
	ID           int
	Name         string
	Description  string
	ImagePath    string
	Price        int
	CreatedAt    string
	CommentCount int
	Comments     []CommentWriter
}

// CommentWriter Model
type CommentWriter struct {
	Content string
	Writer  string
}

func getProduct(pid int) Product {
	pv := ProductDB[pid]
	p := Product{}
	p.ID = pv.ID
	p.Name = pv.Name
	p.Description = pv.Description
	p.ImagePath = pv.ImagePath
	p.Price = pv.Price
	p.CreatedAt = pv.CreatedAt
	return p
	var count int
	row := db.QueryRow("SELECT * FROM products WHERE id = ? LIMIT 1", pid)
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &p.CreatedAt, &count)
	if err != nil {
		panic(err.Error())
	}

	return p
}

func getProductsWithCommentsAt(page int) []ProductWithComments {
	// select 50 products with offset page*50
    s := make([]ProductWithComments, 50)
    copy(s, ProductDB[(199-page)*50+1: (200-page)*50+1])
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
	fmt.Println(len(s))
	return s
	products := make([]ProductWithComments, 0, 50)
	//products := []ProductWithComments{}
	//rows, err := db.Query("SELECT * FROM products ORDER BY id DESC LIMIT 50 OFFSET ?", page*50)
	//rows, err := db.Query("SELECT * FROM products WHERE id BETWEEN ? AND ? ORDER BY id DESC", (199-page)*50+1, (200-page)*50)
	rows, err := db.Query("SELECT * FROM products WHERE id BETWEEN ? AND ? ORDER BY id DESC", (199-page)*50+1, (200-page)*50)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		var cnt int
		p := ProductWithComments{}
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &p.CreatedAt, &cnt)
/*
		// select comment count for the product
		var cnt int
		cnterr := db.QueryRow("SELECT count(*) as count FROM comments WHERE product_id = ?", p.ID).Scan(&cnt)
		if cnterr != nil {
			cnt = 0
		}
		p.CommentCount = cnt

		if cnt > 0 {
			// select 5 comments and its writer for the product
			var cWriters []CommentWriter

			subrows, suberr := db.Query("SELECT c.content, u.name  FROM comments as c INNER JOIN users as u "+
				"ON c.user_id = u.id WHERE c.product_id = ? ORDER BY c.id DESC LIMIT 5", p.ID)
			if suberr != nil {
				subrows = nil
			}

			defer subrows.Close()
			for subrows.Next() {
				//var i int
				//var s string
				var cw CommentWriter
				//subrows.Scan(&i, &i, &i, &cw.Content, &s, &i, &cw.Writer, &s, &s, &s)
				subrows.Scan(&cw.Content, &cw.Writer)
				cWriters = append(cWriters, cw)
			}

			p.Comments = cWriters
		}

		products = append(products, p)
*/
		// select comment count for the product
/*
		cnterr := db.QueryRow("SELECT count(1) as count FROM comments WHERE product_id = ?", p.ID).Scan(&cnt)
		if cnterr != nil {
			cnt = 0
		}
*/
		p.CommentCount = cnt
		if cnt > 0 {
			// select 5 comments and its writer for the product
			var cWriters []CommentWriter

			subrows, suberr := db.Query("SELECT c.content, u.name  FROM comments as c INNER JOIN users as u "+
				"ON c.user_id = u.id WHERE c.product_id = ? ORDER BY c.id DESC LIMIT 5", p.ID)
			//subrows, suberr := db.Query("SELECT sub.content, sub.name, id FROM products INNER JOIN (SELECT c.content as content, u.name as name, c.product_id as product_id FROM comments as c INNER JOIN users as u "+
			//	"ON c.user_id = u.id ORDER BY c.id DESC LIMIT 5) as sub ON sub.product_id = id WHERE id BETWEEN ? AND ?", (199-page)*50+1, (200-page)*50)
			if suberr != nil {
				subrows = nil
			}

			defer subrows.Close()
			for subrows.Next() {
				//var i int
				//var s string
				var cw CommentWriter
				//subrows.Scan(&i, &i, &i, &cw.Content, &s, &i, &cw.Writer, &s, &s, &s)
				subrows.Scan(&cw.Content, &cw.Writer)
				cWriters = append(cWriters, cw)
			}

			p.Comments = cWriters
		}

		products = append(products, p)
	}

	return products
}

func (p *Product) isBought(uid int) bool {
	var count int
	log.Print(uid)
	log.Print(p.ID)
	err := db.QueryRow(
		"SELECT count(1) as count FROM histories WHERE product_id = ? AND user_id = ?",
		p.ID, uid,
	).Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	return count > 0
}
