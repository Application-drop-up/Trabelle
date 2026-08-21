package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/Application-drop-up/Travellle/internal/db"
	"github.com/Application-drop-up/Travellle/internal/domain/countryguide"
)

func main() {
	conn, err := db.NewConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close()

	if err := db.RunMigrations(conn, "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	guides, err := countryGuideSeeds()
	if err != nil {
		log.Fatalf("invalid seed data: %v", err)
	}

	ctx := context.Background()
	for _, guide := range guides {
		if err := upsertCountryGuide(ctx, conn, guide); err != nil {
			log.Fatalf("failed to seed country guide %s: %v", guide.CountryCode, err)
		}
		log.Printf("Seeded country guide: %s (%s)", guide.CountryName, guide.CountryCode)
	}
}

// upsertCountryGuide inserts or updates a country_guides row and replaces
// its items wholesale. Items have no independent identity worth preserving
// across re-seeds, so a delete-then-insert is simpler than diffing.
func upsertCountryGuide(ctx context.Context, conn *sql.DB, guide *countryguide.CountryGuide) error {
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var guideID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO country_guides (country_code, country_name)
		VALUES ($1, $2)
		ON CONFLICT (country_code) DO UPDATE SET country_name = EXCLUDED.country_name, updated_at = NOW()
		RETURNING id`,
		guide.CountryCode, guide.CountryName).Scan(&guideID)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM country_guide_items WHERE country_guide_id = $1`, guideID); err != nil {
		return err
	}

	for index, item := range guide.Items {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO country_guide_items (country_guide_id, category, title, description, url, is_mandatory, sort_order)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			guideID, string(item.Category), item.Title, nullIfEmpty(item.Description), nullIfEmpty(item.URL), item.IsMandatory, index,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func nullIfEmpty(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func countryGuideSeeds() ([]*countryguide.CountryGuide, error) {
	definitions := []struct {
		code  string
		name  string
		items []countryguide.Item
	}{
		{
			code: "TH",
			name: "タイ",
			items: []countryguide.Item{
				{
					Category:    countryguide.CategoryEntryCard,
					Title:       "Thailand Digital Arrival Card (TDAC) の事前申請",
					Description: "出発72時間前からオンラインで申請可能。全ての入国者に必須。",
					URL:         "https://tdac.immigration.go.th",
					IsMandatory: true,
				},
				{
					Category:    countryguide.CategorySIMRecommendation,
					Title:       "空港到着ロビーでプリペイドSIMを購入",
					Description: "AIS・dtacのカウンターが到着ロビーにあり、すぐ設定できる。",
				},
				{
					Category:    countryguide.CategoryPackingTip,
					Title:       "電源プラグはBFタイプ",
					Description: "日本のAタイプ機器はそのまま挿せないことが多いので変換アダプターを持参。",
				},
			},
		},
		{
			code: "TW",
			name: "台湾",
			items: []countryguide.Item{
				{
					Category:    countryguide.CategoryEntryCard,
					Title:       "入国カードのオンライン事前登録",
					Description: "内政部移民署のサイトから事前登録可能。未登録でも空港で記入できる。",
					URL:         "https://oa.immigration.gov.tw",
				},
				{
					Category:    countryguide.CategorySIMRecommendation,
					Title:       "桃園空港到着ロビーでSIMを購入",
					Description: "中華電信・台湾大哥大の観光客向けSIMがすぐ購入できる。",
				},
				{
					Category:    countryguide.CategoryPackingTip,
					Title:       "電源プラグは日本と同じAタイプ",
					Description: "変換アダプターは不要。",
				},
			},
		},
		{
			code: "KR",
			name: "韓国",
			items: []countryguide.Item{
				{
					Category:    countryguide.CategoryEntryCard,
					Title:       "K-ETA(電子渡航認証)の要否を出発前に確認",
					Description: "日本国籍者は免除措置が取られている時期があるため、渡航前に最新の要否を確認すること。",
				},
				{
					Category:    countryguide.CategorySIMRecommendation,
					Title:       "仁川空港到着ロビーでSIM/eSIMを受け取り",
					Description: "事前にオンラインで購入し、空港カウンターで受け取るのがスムーズ。",
				},
				{
					Category:    countryguide.CategoryPackingTip,
					Title:       "電源プラグはSE(C/SE)タイプ",
					Description: "変換アダプターが必要。",
				},
			},
		},
		{
			code: "AU",
			name: "オーストラリア",
			items: []countryguide.Item{
				{
					Category:    countryguide.CategoryEntryCard,
					Title:       "ETA(電子渡航認証, サブクラス601)の事前申請",
					Description: "専用アプリまたはウェブサイトから出発前に申請必須。",
					URL:         "https://immi.homeaffairs.gov.au/visas/getting-a-visa/visa-listing/electronic-travel-authority-601",
					IsMandatory: true,
				},
				{
					Category:    countryguide.CategorySIMRecommendation,
					Title:       "空港到着ロビーでプリペイドSIMを購入",
					Description: "Telstra・Optusのカウンターが主要空港にある。",
				},
				{
					Category:    countryguide.CategoryPackingTip,
					Title:       "電源プラグはOタイプ",
					Description: "変換アダプターが必要。",
				},
			},
		},
		{
			code: "CA",
			name: "カナダ",
			items: []countryguide.Item{
				{
					Category:    countryguide.CategoryEntryCard,
					Title:       "eTA(電子渡航認証)の事前申請",
					Description: "航空機で入国する場合は出発前の取得が必須。",
					URL:         "https://www.canada.ca/en/immigration-refugees-citizenship/services/visit-canada/eta.html",
					IsMandatory: true,
				},
				{
					Category:    countryguide.CategorySIMRecommendation,
					Title:       "事前にeSIMを購入しておくのがおすすめ",
					Description: "主要空港にもRogers・Bell・Teluxのカウンターはあるが割高な場合がある。",
				},
				{
					Category:    countryguide.CategoryPackingTip,
					Title:       "電源プラグは日本と同じAタイプ",
					Description: "電圧は120Vで日本の100Vと近く、変換アダプターは基本的に不要。",
				},
			},
		},
	}

	guides := make([]*countryguide.CountryGuide, 0, len(definitions))
	for _, definition := range definitions {
		guide, err := countryguide.NewCountryGuide(definition.code, definition.name, definition.items)
		if err != nil {
			return nil, err
		}
		guides = append(guides, guide)
	}
	return guides, nil
}
