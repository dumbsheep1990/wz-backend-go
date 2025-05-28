package assembly

import (
	"database/sql"

	"github.com/yourusername/wz-backend-go/internal/application/community/service"
	"github.com/yourusername/wz-backend-go/internal/infrastructure/messaging/eventbus"
	"github.com/yourusername/wz-backend-go/internal/infrastructure/persistence/community"
	"github.com/yourusername/wz-backend-go/internal/infrastructure/persistence/database"
	"github.com/yourusername/wz-backend-go/internal/infrastructure/persistence/user"
	"github.com/yourusername/wz-backend-go/internal/interfaces/grpc/community"
)

// CommunityServiceAssembly assembles all the components for the community service
type CommunityServiceAssembly struct {
	GRPCAdapter *community.CommunityGRPCAdapter
}

// NewCommunityServiceAssembly creates a new CommunityServiceAssembly
func NewCommunityServiceAssembly(db *sql.DB, transactionManager database.TransactionManager) *CommunityServiceAssembly {
	// Create event bus
	eventBus := eventbus.NewEventBus()

	// Create unit of work
	unitOfWork := database.NewUnitOfWork(transactionManager)

	// Create repositories
	communityRepo := community.NewCommunityRepository(db)
	userRepo := user.NewUserRepository(db)

	// Create application service
	communityAppService := service.NewCommunityApplicationService(
		communityRepo,
		userRepo,
		eventBus,
		unitOfWork,
	)

	// Create gRPC adapter
	grpcAdapter := community.NewCommunityGRPCAdapter(communityAppService)

	return &CommunityServiceAssembly{
		GRPCAdapter: grpcAdapter,
	}
}
