import { test, expect } from '@playwright/test';

test.describe('API End-to-End Tests', () => {
  test('Items API CRUD operations', async ({ request }) => {
    // Create an item
    const createResponse = await request.post('/api/items', {
      data: {
        name: 'E2E Test Item',
        price: 999,
        stock: 50,
        code: 'E2E001',
        description: 'Created by E2E test'
      }
    });
    expect(createResponse.ok()).toBeTruthy();
    const createdItem = await createResponse.json();
    expect(createdItem.item.name).toBe('E2E Test Item');

    const itemId = createdItem.item.id;

    // Read the item
    const getResponse = await request.get(`/api/items/${itemId}`);
    expect(getResponse.ok()).toBeTruthy();
    const retrievedItem = await getResponse.json();
    expect(retrievedItem.item.price).toBe(999);

    // Update the item
    const updateResponse = await request.put(`/api/items/${itemId}`, {
      data: {
        name: 'E2E Updated Item',
        price: 1999,
        stock: 100,
        code: 'E2E001',
        description: 'Updated by E2E test'
      }
    });
    expect(updateResponse.ok()).toBeTruthy();
    const updatedItem = await updateResponse.json();
    expect(updatedItem.item.name).toBe('E2E Updated Item');
    expect(updatedItem.item.price).toBe(1999);

    // List items
    const listResponse = await request.get('/api/items');
    expect(listResponse.ok()).toBeTruthy();
    const itemsList = await listResponse.json();
    expect(Array.isArray(itemsList.items)).toBeTruthy();
    const foundItem = itemsList.items.find((i: any) => i.id === itemId);
    expect(foundItem).toBeDefined();

    // Delete the item
    const deleteResponse = await request.delete(`/api/items/${itemId}`);
    expect(deleteResponse.ok()).toBeTruthy();

    // Verify deletion
    const verifyResponse = await request.get(`/api/items/${itemId}`);
    expect(verifyResponse.status()).toBe(404);
  });

  test('Staff API CRUD operations', async ({ request }) => {
    // Create staff
    const createResponse = await request.post('/api/staff', {
      data: {
        name: 'E2E Test Staff',
        storeId: 1,
        password: 'testpass123'
      }
    });
    expect(createResponse.ok()).toBeTruthy();
    const createdStaff = await createResponse.json();
    expect(createdStaff.staff.name).toBe('E2E Test Staff');

    const staffId = createdStaff.staff.id;

    // Read staff
    const getResponse = await request.get(`/api/staff/${staffId}`);
    expect(getResponse.ok()).toBeTruthy();
    const retrievedStaff = await getResponse.json();
    expect(retrievedStaff.staff.name).toBe('E2E Test Staff');

    // Update staff
    const updateResponse = await request.put(`/api/staff/${staffId}`, {
      data: {
        name: 'E2E Updated Staff',
        storeId: 1,
        password: 'updatedpass456'
      }
    });
    expect(updateResponse.ok()).toBeTruthy();
    const updatedStaff = await updateResponse.json();
    expect(updatedStaff.staff.name).toBe('E2E Updated Staff');

    // List staff
    const listResponse = await request.get('/api/staff');
    expect(listResponse.ok()).toBeTruthy();
    const staffList = await listResponse.json();
    expect(Array.isArray(staffList.staff)).toBeTruthy();

    // Delete staff
    const deleteResponse = await request.delete(`/api/staff/${staffId}`);
    expect(deleteResponse.ok()).toBeTruthy();

    // Verify deletion
    const verifyResponse = await request.get(`/api/staff/${staffId}`);
    expect(verifyResponse.status()).toBe(404);
  });

  test('Store API CRUD operations', async ({ request }) => {
    // Create store
    const createResponse = await request.post('/api/stores', {
      data: {
        name: 'E2E Test Store',
        code: 'E2ESTORE'
      }
    });
    expect(createResponse.ok()).toBeTruthy();
    const createdStore = await createResponse.json();
    expect(createdStore.store.name).toBe('E2E Test Store');

    const storeId = createdStore.store.id;

    // Read store
    const getResponse = await request.get(`/api/stores/${storeId}`);
    expect(getResponse.ok()).toBeTruthy();
    const retrievedStore = await getResponse.json();
    expect(retrievedStore.store.code).toBe('E2ESTORE');

    // Update store
    const updateResponse = await request.put(`/api/stores/${storeId}`, {
      data: {
        name: 'E2E Updated Store',
        code: 'E2EUPDATE'
      }
    });
    expect(updateResponse.ok()).toBeTruthy();
    const updatedStore = await updateResponse.json();
    expect(updatedStore.store.name).toBe('E2E Updated Store');

    // List stores
    const listResponse = await request.get('/api/stores');
    expect(listResponse.ok()).toBeTruthy();
    const storesList = await listResponse.json();
    expect(Array.isArray(storesList.stores)).toBeTruthy();

    // Delete store
    const deleteResponse = await request.delete(`/api/stores/${storeId}`);
    expect(deleteResponse.ok()).toBeTruthy();

    // Verify deletion
    const verifyResponse = await request.get(`/api/stores/${storeId}`);
    expect(verifyResponse.status()).toBe(404);
  });

  test('Sales API operations', async ({ request }) => {
    // First, create test items for the sale
    const item1Response = await request.post('/api/items', {
      data: {
        name: 'Sale Test Item 1',
        price: 100,
        stock: 100,
        code: 'SALETEST1'
      }
    });
    const item1 = await item1Response.json();

    const item2Response = await request.post('/api/items', {
      data: {
        name: 'Sale Test Item 2',
        price: 200,
        stock: 100,
        code: 'SALETEST2'
      }
    });
    const item2 = await item2Response.json();

    // Create a sale
    const saleResponse = await request.post('/api/sales', {
      data: {
        storeId: 1,
        staffId: 1,
        totalPrice: 500,
        deposit: 600,
        items: [
          { itemId: item1.item.id, price: 100, quantity: 2 },
          { itemId: item2.item.id, price: 200, quantity: 1 }
        ]
      }
    });
    expect(saleResponse.ok()).toBeTruthy();
    const createdSale = await saleResponse.json();
    expect(createdSale.sale.totalPrice).toBe(500);

    const saleId = createdSale.sale.id;

    // Read the sale
    const getResponse = await request.get(`/api/sales/${saleId}`);
    expect(getResponse.ok()).toBeTruthy();
    const retrievedSale = await getResponse.json();
    expect(retrievedSale.sale.deposit).toBe(600);
    expect(Array.isArray(retrievedSale.sale.items)).toBeTruthy();
    expect(retrievedSale.sale.items.length).toBe(2);

    // List sales
    const listResponse = await request.get('/api/sales');
    expect(listResponse.ok()).toBeTruthy();
    const salesList = await listResponse.json();
    expect(Array.isArray(salesList.sales)).toBeTruthy();

    // Verify stock reduction
    const item1AfterSale = await request.get(`/api/items/${item1.item.id}`);
    const item1Data = await item1AfterSale.json();
    expect(item1Data.item.stock).toBe(98); // 100 - 2

    const item2AfterSale = await request.get(`/api/items/${item2.item.id}`);
    const item2Data = await item2AfterSale.json();
    expect(item2Data.item.stock).toBe(99); // 100 - 1

    // Clean up test items
    await request.delete(`/api/items/${item1.item.id}`);
    await request.delete(`/api/items/${item2.item.id}`);
  });

  test('Error handling - Invalid requests', async ({ request }) => {
    // Try to get non-existent item
    const response404 = await request.get('/api/items/999999');
    expect(response404.status()).toBe(404);

    // Try to create item with invalid data
    const responseBadRequest = await request.post('/api/items', {
      data: {
        name: '', // Empty name should fail validation
        price: -100, // Negative price should fail
        stock: -10 // Negative stock should fail
      }
    });
    expect(responseBadRequest.status()).toBe(400);

    // Try to update non-existent item
    const responseNotFound = await request.put('/api/items/999999', {
      data: {
        name: 'Test',
        price: 100,
        stock: 10
      }
    });
    expect(responseNotFound.status()).toBe(404);

    // Try to delete non-existent item
    const responseDeleteNotFound = await request.delete('/api/items/999999');
    expect(responseDeleteNotFound.status()).toBe(404);
  });
});