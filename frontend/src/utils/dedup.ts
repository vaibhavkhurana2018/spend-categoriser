import { db } from '../db';
import type { TransactionRaw, Transaction } from '../types';

export async function deduplicateTransactions(
  incoming: TransactionRaw[],
  billFileName: string
): Promise<{ newTransactions: Transaction[]; duplicateCount: number }> {
  const existingHashes = new Set(
    (await db.transactions.toArray()).map((t) => t.hash)
  );

  const now = Date.now();
  const newTransactions: Transaction[] = [];
  let duplicateCount = 0;

  for (const tx of incoming) {
    if (existingHashes.has(tx.hash)) {
      duplicateCount++;
      continue;
    }
    existingHashes.add(tx.hash);
    newTransactions.push({
      ...tx,
      billFileName,
      addedAt: now,
    });
  }

  return { newTransactions, duplicateCount };
}
