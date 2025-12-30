package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateAuction_AutoClose(t *testing.T) {
	// Configurar intervalo curto para o teste
	os.Setenv("AUCTION_INTERVAL", "2s")
	defer os.Unsetenv("AUCTION_INTERVAL")

	// Conectar ao MongoDB de teste
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/auctions?authSource=admin"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("auctions_test")
	collection := database.Collection("auctions")

	// Limpar coleção antes do teste
	collection.Drop(ctx)

	// Criar repositório
	auctionRepo := &AuctionRepository{
		Collection: collection,
	}

	// Criar entidade de leilão
	auction, internalErr := auction_entity.CreateAuction(
		"Test Product",
		"Electronics",
		"Test product for auction auto-close",
		auction_entity.New,
	)
	if internalErr != nil {
		t.Fatalf("Failed to create auction entity: %v", internalErr)
	}

	// Criar leilão no banco (inicia a goroutine de fechamento)
	internalErr = auctionRepo.CreateAuction(ctx, auction)
	if internalErr != nil {
		t.Fatalf("Failed to create auction in database: %v", internalErr)
	}

	// Verificar que o leilão foi criado com status Active
	var createdAuction AuctionEntityMongo
	err = collection.FindOne(ctx, bson.M{"_id": auction.Id}).Decode(&createdAuction)
	if err != nil {
		t.Fatalf("Failed to find created auction: %v", err)
	}

	if createdAuction.Status != auction_entity.Active {
		t.Errorf("Expected auction status to be Active, got %v", createdAuction.Status)
	}

	// Aguardar o intervalo + margem de segurança (2s + 1s = 3s)
	time.Sleep(3 * time.Second)

	// Verificar se o status foi atualizado para Completed
	var updatedAuction AuctionEntityMongo
	err = collection.FindOne(ctx, bson.M{"_id": auction.Id}).Decode(&updatedAuction)
	if err != nil {
		t.Fatalf("Failed to find auction after interval: %v", err)
	}

	if updatedAuction.Status != auction_entity.Completed {
		t.Errorf("Expected auction status to be Completed after interval, got %v", updatedAuction.Status)
	}

	t.Log("✅ Auction was automatically closed after the configured interval")

	// Limpar após o teste
	collection.Drop(ctx)
}

func TestGetAuctionInterval(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected time.Duration
	}{
		{
			name:     "Valid duration from env",
			envValue: "30s",
			expected: 30 * time.Second,
		},
		{
			name:     "Invalid duration - fallback to default",
			envValue: "invalid",
			expected: 5 * time.Minute,
		},
		{
			name:     "Empty env - fallback to default",
			envValue: "",
			expected: 5 * time.Minute,
		},
		{
			name:     "Complex duration",
			envValue: "1m30s",
			expected: 90 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("AUCTION_INTERVAL", tt.envValue)
				defer os.Unsetenv("AUCTION_INTERVAL")
			} else {
				os.Unsetenv("AUCTION_INTERVAL")
			}

			result := getAuctionInterval()

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
