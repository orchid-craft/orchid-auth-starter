package repository

import (
	"context"

	authModel "github.com/orchid-craft/orchid-auth-starter/modules/auth/domain/models"
)

type SessionRepositoryInterface interface {
	Save(ctx context.Context, session *authModel.RefreshTokenSession) error
	// FindByID mengembalikan SALINAN sesi, bukan pointer ke data internal repository,
	// sehingga pemanggil tidak bisa membaca/memutasi state bersama tanpa lock.
	FindByID(ctx context.Context, id string) (*authModel.RefreshTokenSession, error)
	// Consume menandai sesi sebagai terpakai (revoked) secara ATOMIK untuk Refresh Token Rotation.
	//   - sesi tidak ada      -> domain.ErrTokenInvalid
	//   - sesi sudah revoked  -> domain.ErrTokenReused (reuse attempt atau request paralel)
	//   - sesi aktif          -> Revoked di-set true, salinan sesi dikembalikan
	// Untuk satu token yang sama, hanya SATU pemanggil konkuren yang bisa berhasil.
	// Implementasi database harus setara dengan: UPDATE ... SET revoked = true WHERE id = $1 AND revoked = false.
	Consume(ctx context.Context, id string) (*authModel.RefreshTokenSession, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllForUser(ctx context.Context, userID string) error
}
