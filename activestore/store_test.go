package activestore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/entitystore"
	_ "modernc.org/sqlite"
)

func initStore(t *testing.T, name string) entitystore.StoreInterface {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+name+"?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                 db,
		EntityTableName:    "entity",
		AttributeTableName: "attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func TestNew_NilStore(t *testing.T) {
	s, err := New(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil store")
	}
	if s != nil {
		t.Fatal("expected nil ActiveStoreInterface")
	}
}

func TestNewEntity_FluentSave(t *testing.T) {
	ctx := context.Background()
	store := initStore(t, "activestore_fluent_save.db")
	active, err := New(ctx, store)
	if err != nil {
		t.Fatal(err)
	}

	product := active.EntityCreate("product").
		SetString("name", "Laptop").
		SetFloat("price", 1299.99).
		SetInt("stock", 50)

	if _, err := product.Save(); err != nil {
		t.Fatal(err)
	}

	id := product.GetEntity().ID()
	name, exists, err := product.GetString("name")
	if err != nil || !exists || name != "Laptop" {
		t.Fatalf("expected name=Laptop, got %q exists=%v err=%v", name, exists, err)
	}

	found, err := active.EntityFindByID(id)
	if err != nil {
		t.Fatal(err)
	}
	if found.GetEntity().GetType() != "product" {
		t.Fatal("expected type product")
	}
}

func TestFindByID_AndImmediateWrites(t *testing.T) {
	ctx := context.Background()
	store := initStore(t, "activestore_find.db")
	active, err := New(ctx, store)
	if err != nil {
		t.Fatal(err)
	}

	entity, err := store.EntityCreateWithType(ctx, "product")
	if err != nil {
		t.Fatal(err)
	}

	product, err := active.EntityFindByID(entity.ID())
	if err != nil {
		t.Fatal(err)
	}

	// persisted entity: writes go to the store immediately
	product.SetString("name", "Phone")
	if err := product.Err(); err != nil {
		t.Fatal(err)
	}

	name, exists, err := store.AttributeGetString(ctx, entity.ID(), "name")
	if err != nil || !exists || name != "Phone" {
		t.Fatalf("expected name=Phone, got %q exists=%v err=%v", name, exists, err)
	}
}

func TestListAndCount(t *testing.T) {
	ctx := context.Background()
	store := initStore(t, "activestore_list.db")
	active, err := New(ctx, store)
	if err != nil {
		t.Fatal(err)
	}

	for range 3 {
		if _, err := active.EntityCreate("tag").Save(); err != nil {
			t.Fatal(err)
		}
	}

	query := entitystore.EntityQuery().WithEntityType("tag")
	count, err := active.EntityCount(query)
	if err != nil || count != 3 {
		t.Fatalf("expected count=3, got %v err=%v", count, err)
	}

	list, err := active.EntityList(query)
	if err != nil || len(list) != 3 {
		t.Fatalf("expected 3 items, got %v err=%v", len(list), err)
	}
}

func TestTrashAndDelete(t *testing.T) {
	ctx := context.Background()
	store := initStore(t, "activestore_trash.db")
	active, err := New(ctx, store)
	if err != nil {
		t.Fatal(err)
	}

	product := active.EntityCreate("product")
	if _, err := product.Save(); err != nil {
		t.Fatal(err)
	}
	id := product.GetEntity().ID()

	trashed, err := product.Trash()
	if err != nil || !trashed {
		t.Fatalf("expected trashed=true, got %v err=%v", trashed, err)
	}
	_ = id

	other := active.EntityCreate("product")
	if _, err := other.Save(); err != nil {
		t.Fatal(err)
	}
	deleted, err := other.Delete()
	if err != nil || !deleted {
		t.Fatalf("expected deleted=true, got %v err=%v", deleted, err)
	}
}

func TestPrefetch(t *testing.T) {
	ctx := context.Background()
	store := initStore(t, "activestore_prefetch.db")
	active, err := New(ctx, store)
	if err != nil {
		t.Fatal(err)
	}

	product := active.EntityCreate("product").
		SetString("name", "Laptop").
		SetInt("stock", 50)
	if _, err := product.Save(); err != nil {
		t.Fatal(err)
	}

	found, err := active.EntityFindByID(product.GetEntity().ID())
	if err != nil {
		t.Fatal(err)
	}
	if err := found.Prefetch(); err != nil {
		t.Fatal(err)
	}

	name, exists, err := found.GetString("name")
	if err != nil || !exists || name != "Laptop" {
		t.Fatalf("expected name=Laptop, got %q exists=%v err=%v", name, exists, err)
	}

	stock, exists, _ := found.GetString("stock")
	if !exists || stock != "50" {
		t.Fatalf("expected stock=50, got %q exists=%v", stock, exists)
	}

	missing, exists, _ := found.GetString("nonexistent")
	if exists || missing != "" {
		t.Fatalf("expected missing key, got %q exists=%v", missing, exists)
	}

	// setters keep the cache in sync
	found.SetString("name", "Tablet")
	if err := found.Err(); err != nil {
		t.Fatal(err)
	}
	name, _, _ = found.GetString("name")
	if name != "Tablet" {
		t.Fatalf("expected cached name=Tablet, got %q", name)
	}
}

func TestWrapEntity_Nil(t *testing.T) {
	active, err := New(context.Background(), initStore(t, "activestore_wrap.db"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := active.EntityWrap(nil); err == nil {
		t.Fatal("expected error for nil entity")
	}
}
