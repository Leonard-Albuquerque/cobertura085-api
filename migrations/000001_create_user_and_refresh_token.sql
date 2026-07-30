-- =========================================================
-- Camada de autenticação (JWT) para o schema de Stores
-- Regra de negócio: 1 usuário está associado a exatamente 1 Store (1:1)
-- =========================================================

-- Usuário: guarda os dados de login, já vinculado à loja dele
CREATE TABLE IF NOT EXISTS "User" (
	"id" text PRIMARY KEY,
	"storeId" text NOT NULL,
	"email" text NOT NULL,
	"passwordHash" text NOT NULL,
	"name" text,
	"isActive" boolean DEFAULT true NOT NULL,
	"lastLoginAt" timestamp,
	"createdAt" timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updatedAt" timestamp NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "User_pkey" ON "User" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "User_email_key" ON "User" ("email");
-- Garante a regra 1:1 no nível do banco: uma loja só pode ter um usuário
CREATE UNIQUE INDEX IF NOT EXISTS "User_storeId_key" ON "User" ("storeId");

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'User_storeId_fkey'
    ) THEN
        ALTER TABLE "User" ADD CONSTRAINT "User_storeId_fkey"
            FOREIGN KEY ("storeId") REFERENCES "Store"("id") ON DELETE CASCADE ON UPDATE CASCADE;
    END IF;
END $$;

-- Refresh tokens: permite revogar sessões sem esperar o access token
-- expirar, e rastrear de onde veio o login
CREATE TABLE IF NOT EXISTS "RefreshToken" (
	"id" text PRIMARY KEY,
	"userId" text NOT NULL,
	"tokenHash" text NOT NULL, -- nunca guardar o token em texto puro
	"expiresAt" timestamp NOT NULL,
	"revokedAt" timestamp,
	"userAgent" text,
	"ipHash" varchar(64),
	"createdAt" timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "RefreshToken_pkey" ON "RefreshToken" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "RefreshToken_tokenHash_key" ON "RefreshToken" ("tokenHash");
CREATE INDEX IF NOT EXISTS "RefreshToken_userId_idx" ON "RefreshToken" ("userId");

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'RefreshToken_userId_fkey'
    ) THEN
        ALTER TABLE "RefreshToken" ADD CONSTRAINT "RefreshToken_userId_fkey"
            FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;
    END IF;
END $$;
