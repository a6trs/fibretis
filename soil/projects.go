package soil

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type Project struct {
	ID          int
	Title       string
	Desc        string
	Author      int
	State       int
	TitleColour string
	BannerImg   string
	BannerType  int
	CreatedAt   time.Time
}

func init_Project() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		desc TEXT,
		author INTEGER,
		state INTEGER,
		title_clr VARCHAR(7),
		banner_img TEXT,
		banner_type INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(author) REFERENCES accounts(id)
	)`)
	return err
}

const (
	KEY_Project_ID = iota
	KEY_Project_State
)

const (
	Project_StUnsaved = iota
	Project_StPurposed
	Project_StSeeded
	Project_StRooting
	Project_StSprouts
	Project_StLawn
	Project_StWood
	Project_StJungle
	Project_StForest
)

const (
	BI_Pattern = iota
	BI_Cover
)

func StateStyles(state int) (string, string) {
	switch state {
	case Project_StUnsaved:
		return "#999999", "Unsaved"
	case Project_StPurposed:
		return "#0dcfc7", "Purposed"
	case Project_StSeeded:
		return "#9f8e0e", "Seeded"
	case Project_StRooting:
		return "#e2d904", "Rooting"
	case Project_StSprouts:
		return "#76e331", "Sprouts"
	case Project_StLawn:
		return "#01dd63", "Lawn"
	case Project_StWood:
		return "#16bc08", "Wood"
	case Project_StJungle:
		return "#049b27", "Jungle"
	case Project_StForest:
		return "#046009", "Forest"
	default:
		return "#999999", "Unknown"
	}
}

// Usage: class='banner <COBT(BT)>'
func ClassOfBannerType(bitype int) string {
	switch bitype {
	case BI_Pattern:
		return "bi-pattern"
	case BI_Cover:
		return "bi-cover"
	default:
		return ""
	}
}

func (this *Project) Find(key int) int {
	result := -1
	var row *sql.Row
	switch key {
	case KEY_Project_ID:
		row = db.QueryRow(`SELECT id FROM projects WHERE id = ?`, this.ID)
	case KEY_Project_State:
		row = db.QueryRow(`SELECT id FROM projects WHERE state = ?`, this.State)
	default:
		return -1
	}
	err := row.Scan(&result)
	if err == nil {
		return result
	} else {
		return -1
	}
}

func (this *Project) Load(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		return ErrRowNotFound
	}
	row := db.QueryRow(`SELECT * FROM projects WHERE id = ?`, this.ID)
	return row.Scan(&this.ID, &this.Title, &this.Desc, &this.Author, &this.State, &this.TitleColour, &this.BannerImg, &this.BannerType, &this.CreatedAt)
}

func (this *Project) Save(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		_, err := db.Exec(`INSERT INTO projects (state) VALUES (?)`, Project_StUnsaved)
		if err != nil {
			return err
		}
		state := this.State
		this.State = Project_StUnsaved
		this.ID = this.Find(KEY_Project_State)
		this.State = state
	}
	_, err := db.Exec(`UPDATE projects SET title = ?, desc = ?, author = ?, state = ?, title_clr = ?, banner_img = ?, banner_type = ? WHERE id = ?`, this.Title, this.Desc, this.Author, this.State, this.TitleColour, this.BannerImg, this.BannerType, this.ID)
	return err
}

func NumberOfProjects() int {
	var n int
	row := db.QueryRow(`SELECT COUNT(*) FROM projects`)
	if row.Scan(&n) == nil {
		return n
	} else {
		return -1
	}
}

func RecommendProjects(prjid int) []int {
	var list []int
	// TODO: Improve this algorithm whenever possible
	// Here we just retrieve all the people whose sight level for this project
	//   is not zero and find what else they stared at (or you can say starred)
	rs1, err := db.Query(`SELECT account FROM sights_projects WHERE target = ? AND level <> 0`, prjid)
	if err != nil {
		return nil
	}
	defer rs1.Close()
	gazers := []string{}
	for rs1.Next() {
		var a int
		if rs1.Scan(&a) == nil {
			gazers = append(gazers, strconv.Itoa(a))
		}
	}
	// stackoverflow.com/q/1503959
	rs2, err := db.Query(`SELECT target FROM sights_projects WHERE account IN (`+strings.Join(gazers, ",")+`) AND target <> ? GROUP BY target ORDER BY count(*) DESC LIMIT 3`, prjid)
	if err != nil {
		return nil
	}
	defer rs2.Close()
	for rs2.Next() {
		var a int
		if rs2.Scan(&a) == nil {
			list = append(list, a)
		}
	}
	return list
}
