import Dexie, { type Table } from 'dexie';
import type { Transaction, Bill } from '../types';

class SpendDB extends Dexie {
  transactions!: Table<Transaction, number>;
  bills!: Table<Bill, number>;

  constructor() {
    super('SpendCategoriser');
    this.version(1).stores({
      transactions: '++id, hash, date, category, subCategory, billFileName',
      bills: '++id, fileName, uploadedAt',
    });
  }
}

export const db = new SpendDB();
