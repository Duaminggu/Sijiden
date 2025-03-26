package seeder

import (
	"context"
	"fmt"

	"github.com/duaminggu/sijiden/ent"
	"github.com/duaminggu/sijiden/ent/role"
	"github.com/duaminggu/sijiden/ent/user"
	"github.com/duaminggu/sijiden/ent/userrole"
)

func Seed(client *ent.Client) error {
	ctx := context.Background()

	// Seed Role
	adminRole, err := client.Role.Query().Where(role.NameEQ("admin")).Only(ctx)
	if ent.IsNotFound(err) {
		adminRole, err = client.Role.
			Create().
			SetName("admin").
			SetDescription("Administrator Role").
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed seeding admin role: %w", err)
		}
		fmt.Println("✅ Role 'admin' created.")
	} else {
		fmt.Println("ℹ️  Role 'admin' already exists.")
	}

	// Seed User
	adminUser, err := client.User.Query().Where(user.UsernameEQ("superadmin")).Only(ctx)
	if ent.IsNotFound(err) {
		adminUser, err = client.User.
			Create().
			SetUsername("superadmin").
			SetEmail("admin@mail.com").
			SetPassword("ZXasqw12!@").
			SetEmailVerified(true).
			SetPhoneVerified(true).
			SetFirstName("Super").
			SetLastName("Admin").
			SetPhoneNumber("0800000000").
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed seeding admin user: %w", err)
		}
		fmt.Println("✅ User 'superadmin' created.")
	} else {
		fmt.Println("ℹ️  User 'superadmin' already exists.")
	}

	// Assign Role to User (UserRole)
	_, err = client.UserRole.
		Query().
		Where(
			userrole.HasUserWith(user.IDEQ(adminUser.ID)),
			userrole.HasRoleWith(role.IDEQ(adminRole.ID)),
		).
		Only(ctx)

	if ent.IsNotFound(err) {
		_, err = client.UserRole.
			Create().
			SetUserID(adminUser.ID).
			SetRoleID(adminRole.ID).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed seeding user-role: %w", err)
		}
		fmt.Println("✅ Role 'admin' assigned to 'superadmin'.")
	} else {
		fmt.Println("ℹ️  'superadmin' already has 'admin' role.")
	}

	return nil
}

func SeedIfNeeded(client *ent.Client) error {
	ctx := context.Background()

	// Cek apakah sudah ada data user atau role
	userCount, err := client.User.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to check user count: %w", err)
	}
	roleCount, err := client.Role.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to check role count: %w", err)
	}

	if userCount > 0 || roleCount > 0 {
		fmt.Println("ℹ️  Skipping seed: data already exists.")
		return nil
	}

	fmt.Println("🚀 Running seed...")
	return Seed(client)
}
