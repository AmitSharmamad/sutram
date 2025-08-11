package migrate

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"algo/internal/infra/persistance/repository"
)

// DhanMapper holds the column mapping configuration
type DhanMapper struct {
	ColumnMapping          map[string]string            `json:"column_mapping"`
	EnumMappings           map[string]map[string]string `json:"enum_mappings"`
	SutramColumnsNotInDhan []string                     `json:"sutram_columns_not_in_dhan"`
	DhanColumnsNotInSutram []string                     `json:"dhan_columns_not_in_sutram"`
}

// InstrumentMaster represents the instrument master record
type InstrumentMaster struct {
	ExchangeSecurityID       int32           `db:"exchange_security_id"`
	Exchange                 string          `db:"exchange"`
	Segment                  string          `db:"segment"`
	Symbol                   string          `db:"symbol"`
	Instrument               string          `db:"instrument"`
	InstrumentName           string          `db:"instrument_name"`
	ISIN                     string          `db:"isin"`
	UnderlyingSecurityID     sql.NullInt32   `db:"underlying_security_id"`
	UnderlyingSymbol         sql.NullString  `db:"underlying_symbol"`
	Expiry                   time.Time       `db:"expiry"`
	ExpiryFlag               sql.NullString  `db:"expiry_flag"`
	LotUnits                 int32           `db:"lot_units"`
	TickSize                 float64         `db:"tick_size"`
	StrikePrice              sql.NullFloat64 `db:"strike_price"`
	OptionType               sql.NullString  `db:"option_type"`
	InstrumentType           string          `db:"instrument_type"`
	BracketFlag              string          `db:"bracket_flag"`
	CoverFlag                string          `db:"cover_flag"`
	AsmGsmFlag               string          `db:"asm_gsm_flag"`
	BuySellIndicator         string          `db:"buy_sell_indicator"`
	BuyCoMinMarginPer        sql.NullFloat64 `db:"buy_co_min_margin_per"`
	SellCoMinMarginPer       sql.NullFloat64 `db:"sell_co_min_margin_per"`
	BuyCoSlRangeMaxPerc      sql.NullFloat64 `db:"buy_co_sl_range_max_perc"`
	SellCoSlRangeMinPerc     sql.NullFloat64 `db:"sell_co_sl_range_min_perc"`
	BuyBoMinMarginPer        sql.NullFloat64 `db:"buy_bo_min_margin_per"`
	BuyBoSlRangeMaxPerc      sql.NullFloat64 `db:"buy_bo_sl_range_max_perc"`
	SellBoSlMinRange         sql.NullFloat64 `db:"sell_bo_sl_min_range"`
	BuyBoProfitRangeMaxPerc  sql.NullFloat64 `db:"buy_bo_profit_range_max_perc"`
	SellBoProfitRangeMaxPerc sql.NullFloat64 `db:"sell_bo_profit_range_max_perc"`
	MtfLeverage              sql.NullFloat64 `db:"mtf_leverage"`
	CreatedAt                time.Time       `db:"created_at"`
	UpdatedAt                time.Time       `db:"updated_at"`
}

// InstrumentMasterMigrator handles migration of instrument master data from Excel
type InstrumentMasterMigrator struct {
	db     *sqlx.DB
	mapper *DhanMapper
}

// NewInstrumentMasterMigrator creates a new instrument master migrator
func NewInstrumentMasterMigrator() (*InstrumentMasterMigrator, error) {
	// Load the mapping configuration
	mapper, err := loadDhanMapper()
	if err != nil {
		return nil, fmt.Errorf("failed to load dhan mapper: %w", err)
	}

	// Get ClickHouse connection
	db, err := repository.ProvideClickHouseDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	return &InstrumentMasterMigrator{
		db:     db,
		mapper: mapper,
	}, nil
}

// MigrateFromCSV migrates instrument master data from a CSV file
func (m *InstrumentMasterMigrator) MigrateFromCSV(filePath string) error {
	startTime := time.Now()
	fmt.Printf("Starting migration from CSV file: %s at %s\n", filePath, startTime.Format("2006-01-02 15:04:05"))

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least 2 rows (header + data)")
	}

	// Process the data
	header := records[0]
	dataRows := records[1:]

	fmt.Printf("Found %d header columns: %v\n", len(header), header[:min(10, len(header))])
	fmt.Printf("Found %d data rows to process\n", len(dataRows))

	// Validate required columns are present
	err = m.validateRequiredColumns(header)
	if err != nil {
		return fmt.Errorf("column validation failed: %w", err)
	}

	// Create column mapping from Dhan to Sutram
	columnMap := make(map[int]string)
	mappedColumns := 0
	for i, dhanColumn := range header {
		if sutramColumn, exists := m.mapper.ColumnMapping[dhanColumn]; exists {
			columnMap[i] = sutramColumn
			mappedColumns++
		}
	}
	fmt.Printf("Successfully mapped %d columns\n", mappedColumns)

	// Process data in batches
	batchSize := 1000
	totalProcessed := 0
	totalSkipped := 0
	batchStartTime := time.Now()

	for i := 0; i < len(dataRows); i += batchSize {
		end := i + batchSize
		if end > len(dataRows) {
			end = len(dataRows)
		}

		batch := dataRows[i:end]
		processed, skipped, err := m.processBatch(batch, columnMap, header, i+2) // +2 for 1-indexed and header row
		if err != nil {
			return fmt.Errorf("failed to process batch starting at row %d: %w", i+2, err)
		}

		totalProcessed += processed
		totalSkipped += skipped

		// Calculate batch timing
		batchEndTime := time.Now()
		batchDuration := batchEndTime.Sub(batchStartTime)

		fmt.Printf("Processed %d/%d rows (skipped: %d) - Batch time: %s\n",
			totalProcessed, len(dataRows), totalSkipped, batchDuration.String())

		batchStartTime = time.Now() // Reset for next batch
	}

	fmt.Printf("Migration completed successfully. Processed %d rows, skipped %d rows.\n", totalProcessed, totalSkipped)

	// Calculate and display total time taken
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("Migration finished at %s\n", endTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Total time taken: %s\n", duration.String())

	if totalProcessed > 0 {
		avgTimePerRow := duration / time.Duration(totalProcessed)
		fmt.Printf("Average time per row: %s\n", avgTimePerRow.String())
	}

	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// processBatch processes a batch of rows and inserts them into the database
func (m *InstrumentMasterMigrator) processBatch(batch [][]string, columnMap map[int]string, header []string, startRowNum int) (int, int, error) {
	instruments := make([]*InstrumentMaster, 0, len(batch))
	skipped := 0

	for i, row := range batch {
		rowNum := startRowNum + i
		instrument, err := m.convertRowToInstrument(row, columnMap, header)
		if err != nil {
			fmt.Printf("Warning: Skipping row %d due to conversion error: %v\n", rowNum, err)
			skipped++
			continue
		}
		instruments = append(instruments, instrument)
	}

	if len(instruments) == 0 {
		return 0, skipped, nil // No valid instruments to insert
	}

	err := m.insertInstruments(instruments)
	if err != nil {
		return 0, skipped, err
	}

	return len(instruments), skipped, nil
}

// convertRowToInstrument converts a row of Excel data to an InstrumentMaster struct
func (m *InstrumentMasterMigrator) convertRowToInstrument(row []string, columnMap map[int]string, header []string) (*InstrumentMaster, error) {
	instrument := &InstrumentMaster{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create a map for easier field access
	dhanData := make(map[string]string)
	for i, value := range row {
		if i < len(header) {
			dhanData[header[i]] = value
		}
	}

	// Map each field
	for i, value := range row {
		if i >= len(header) {
			break
		}

		sutramColumn, exists := columnMap[i]
		if !exists {
			continue // Skip unmapped columns
		}

		err := m.setFieldValue(instrument, sutramColumn, value, dhanData)
		if err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", sutramColumn, err)
		}
	}

	return instrument, nil
}

// setFieldValue sets the appropriate field value with type conversion and enum mapping
func (m *InstrumentMasterMigrator) setFieldValue(instrument *InstrumentMaster, fieldName, value string, dhanData map[string]string) error {
	// Trim whitespace
	value = strings.TrimSpace(value)

	switch fieldName {
	case "exchange_security_id":
		if value == "" {
			return fmt.Errorf("exchange_security_id cannot be empty")
		}
		val, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid exchange_security_id: %w", err)
		}
		instrument.ExchangeSecurityID = int32(val)

	case "exchange":
		mappedValue, exists := m.mapper.EnumMappings["exchange"][value]
		if !exists {
			return fmt.Errorf("unknown exchange value: %s", value)
		}
		instrument.Exchange = mappedValue

	case "segment":
		mappedValue, exists := m.mapper.EnumMappings["segment"][value]
		if !exists {
			return fmt.Errorf("unknown segment value: %s", value)
		}
		instrument.Segment = mappedValue

	case "symbol":
		instrument.Symbol = value

	case "instrument":
		mappedValue, exists := m.mapper.EnumMappings["instrument"][value]
		if !exists {
			return fmt.Errorf("unknown instrument value: %s", value)
		}
		instrument.Instrument = mappedValue

	case "instrument_name":
		instrument.InstrumentName = value

	case "isin":
		if value != "" && value != "NA" && value != "null" && value != "NULL" {
			instrument.ISIN = value
		}

	case "underlying_security_id":
		if value != "" && value != "null" && value != "NULL" && value != "0" {
			val, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid underlying_security_id: %w", err)
			}
			instrument.UnderlyingSecurityID = sql.NullInt32{Int32: int32(val), Valid: true}
		}

	case "underlying_symbol":
		if value != "" && value != "null" && value != "NULL" && value != "NA" {
			instrument.UnderlyingSymbol = sql.NullString{String: value, Valid: true}
		}

	case "expiry":
		if value != "" {
			// Try multiple date formats
			formats := []string{
				"2006-01-02 15:04:05",
				"2006-01-02",
				"02-01-2006",
				"02/01/2006",
			}
			var parsedTime time.Time
			var err error
			for _, format := range formats {
				parsedTime, err = time.Parse(format, value)
				if err == nil {
					break
				}
			}
			if err != nil {
				return fmt.Errorf("invalid expiry date format: %s", value)
			}
			instrument.Expiry = parsedTime
		} else {
			// Set default expiry for non-derivative instruments
			instrument.Expiry = time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC)
		}

	case "expiry_flag":
		if value != "" {
			instrument.ExpiryFlag = sql.NullString{String: value, Valid: true}
		}

	case "lot_units":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid lot_units: %w", err)
		}
		instrument.LotUnits = int32(val)

	case "tick_size":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid tick_size: %w", err)
		}
		instrument.TickSize = val

	case "strike_price":
		if value != "" && value != "0" && value != "null" && value != "NULL" && value != "-0.01000" {
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid strike_price: %w", err)
			}
			instrument.StrikePrice = sql.NullFloat64{Float64: val, Valid: true}
		}

	case "option_type":
		if value != "" {
			instrument.OptionType = sql.NullString{String: value, Valid: true}
		}

	case "instrument_type":
		instrument.InstrumentType = value

	case "bracket_flag":
		mappedValue, exists := m.mapper.EnumMappings["bracket_flag"][value]
		if !exists {
			mappedValue = "N" // Default to N
		}
		instrument.BracketFlag = mappedValue

	case "cover_flag":
		mappedValue, exists := m.mapper.EnumMappings["cover_flag"][value]
		if !exists {
			mappedValue = "N" // Default to N
		}
		instrument.CoverFlag = mappedValue

	case "asm_gsm_flag":
		mappedValue, exists := m.mapper.EnumMappings["asm_gsm_flag"][value]
		if !exists {
			mappedValue = "N" // Default to N
		}
		instrument.AsmGsmFlag = mappedValue

	case "buy_sell_indicator":
		if value == "" {
			value = "A" // Default to Both Buy/Sell allowed
		}
		instrument.BuySellIndicator = value

	// Handle nullable float fields
	case "buy_co_min_margin_per", "sell_co_min_margin_per", "buy_co_sl_range_max_perc",
		"sell_co_sl_range_min_perc", "buy_bo_min_margin_per", "buy_bo_sl_range_max_perc",
		"sell_bo_sl_min_range", "buy_bo_profit_range_max_perc", "sell_bo_profit_range_max_perc",
		"mtf_leverage":
		if value != "" && value != "0" && value != "null" && value != "NULL" {
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid %s: %w", fieldName, err)
			}
			nullFloat := sql.NullFloat64{Float64: val, Valid: true}

			switch fieldName {
			case "buy_co_min_margin_per":
				instrument.BuyCoMinMarginPer = nullFloat
			case "sell_co_min_margin_per":
				instrument.SellCoMinMarginPer = nullFloat
			case "buy_co_sl_range_max_perc":
				instrument.BuyCoSlRangeMaxPerc = nullFloat
			case "sell_co_sl_range_min_perc":
				instrument.SellCoSlRangeMinPerc = nullFloat
			case "buy_bo_min_margin_per":
				instrument.BuyBoMinMarginPer = nullFloat
			case "buy_bo_sl_range_max_perc":
				instrument.BuyBoSlRangeMaxPerc = nullFloat
			case "sell_bo_sl_min_range":
				instrument.SellBoSlMinRange = nullFloat
			case "buy_bo_profit_range_max_perc":
				instrument.BuyBoProfitRangeMaxPerc = nullFloat
			case "sell_bo_profit_range_max_perc":
				instrument.SellBoProfitRangeMaxPerc = nullFloat
			case "mtf_leverage":
				instrument.MtfLeverage = nullFloat
			}
		}
	}

	return nil
}

// insertInstruments inserts a batch of instruments into the database
func (m *InstrumentMasterMigrator) insertInstruments(instruments []*InstrumentMaster) error {
	if len(instruments) == 0 {
		return nil
	}

	// Build the bulk insert query
	query := `INSERT INTO sutram.instrument_master (
		exchange_security_id, exchange, segment, symbol, instrument, instrument_name,
		isin, underlying_security_id, underlying_symbol, expiry, expiry_flag,
		lot_units, tick_size, strike_price, option_type, instrument_type,
		bracket_flag, cover_flag, asm_gsm_flag, buy_sell_indicator,
		buy_co_min_margin_per, sell_co_min_margin_per, buy_co_sl_range_max_perc,
		sell_co_sl_range_min_perc, buy_bo_min_margin_per, buy_bo_sl_range_max_perc,
		sell_bo_sl_min_range, buy_bo_profit_range_max_perc, sell_bo_profit_range_max_perc,
		mtf_leverage, created_at, updated_at
	) VALUES `

	// Build values placeholders
	valuePlaceholders := make([]string, len(instruments))
	args := make([]interface{}, 0, len(instruments)*32)

	for i, instrument := range instruments {
		valuePlaceholders[i] = `(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		args = append(args,
			instrument.ExchangeSecurityID,
			instrument.Exchange,
			instrument.Segment,
			instrument.Symbol,
			instrument.Instrument,
			instrument.InstrumentName,
			instrument.ISIN,
			instrument.UnderlyingSecurityID,
			instrument.UnderlyingSymbol,
			instrument.Expiry,
			instrument.ExpiryFlag,
			instrument.LotUnits,
			instrument.TickSize,
			instrument.StrikePrice,
			instrument.OptionType,
			instrument.InstrumentType,
			instrument.BracketFlag,
			instrument.CoverFlag,
			instrument.AsmGsmFlag,
			instrument.BuySellIndicator,
			instrument.BuyCoMinMarginPer,
			instrument.SellCoMinMarginPer,
			instrument.BuyCoSlRangeMaxPerc,
			instrument.SellCoSlRangeMinPerc,
			instrument.BuyBoMinMarginPer,
			instrument.BuyBoSlRangeMaxPerc,
			instrument.SellBoSlMinRange,
			instrument.BuyBoProfitRangeMaxPerc,
			instrument.SellBoProfitRangeMaxPerc,
			instrument.MtfLeverage,
			instrument.CreatedAt,
			instrument.UpdatedAt,
		)
	}

	query += strings.Join(valuePlaceholders, ", ")

	// Execute the query
	_, err := m.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert instruments: %w", err)
	}

	return nil
}

// loadDhanMapper loads the Dhan mapper configuration from JSON file
func loadDhanMapper() (*DhanMapper, error) {
	file, err := os.Open("internal/modules/migrate/dhan_mapper.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open mapper file: %w", err)
	}
	defer file.Close()

	var mapper DhanMapper
	err = json.NewDecoder(file).Decode(&mapper)
	if err != nil {
		return nil, fmt.Errorf("failed to decode mapper JSON: %w", err)
	}

	return &mapper, nil
}

// validateRequiredColumns validates that required columns are present in the Excel header
func (m *InstrumentMasterMigrator) validateRequiredColumns(header []string) error {
	requiredDhanColumns := []string{
		"SECURITY_ID",
		"EXCH_ID",
		"SEGMENT",
		"SYMBOL_NAME",
		"INSTRUMENT",
		"DISPLAY_NAME",
		"LOT_SIZE",
		"TICK_SIZE",
	}

	headerMap := make(map[string]bool)
	for _, col := range header {
		headerMap[col] = true
	}

	var missingColumns []string
	for _, required := range requiredDhanColumns {
		if !headerMap[required] {
			missingColumns = append(missingColumns, required)
		}
	}

	if len(missingColumns) > 0 {
		return fmt.Errorf("missing required columns: %v", missingColumns)
	}

	return nil
}

// TruncateInstrumentMaster truncates the instrument_master table
func (m *InstrumentMasterMigrator) TruncateInstrumentMaster() error {
	_, err := m.db.Exec("TRUNCATE TABLE sutram.instrument_master")
	if err != nil {
		return fmt.Errorf("failed to truncate instrument_master table: %w", err)
	}
	fmt.Println("Truncated instrument_master table")
	return nil
}
