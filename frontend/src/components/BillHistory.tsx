import { useTransactions } from '../hooks/useTransactions';

export default function BillHistory() {
  const { bills } = useTransactions();

  if (bills.length === 0) {
    return (
      <div className="text-center py-20">
        <div className="text-5xl mb-4">📁</div>
        <h2 className="text-xl font-bold text-gray-900">No bills uploaded</h2>
        <p className="text-gray-500 mt-1">Upload credit card statements to see history</p>
      </div>
    );
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">Bill History</h2>
      <p className="text-gray-500 mt-1">{bills.length} bills uploaded</p>

      <div className="mt-6 bg-white rounded-xl border border-gray-200 overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-600">File Name</th>
              <th className="px-4 py-3 text-left font-medium text-gray-600">Uploaded</th>
              <th className="px-4 py-3 text-left font-medium text-gray-600">Period</th>
              <th className="px-4 py-3 text-right font-medium text-gray-600">Amount</th>
              <th className="px-4 py-3 text-right font-medium text-gray-600">Transactions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {bills.map((bill) => (
              <tr key={bill.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 text-gray-900 font-medium">{bill.fileName}</td>
                <td className="px-4 py-3 text-gray-600">{new Date(bill.uploadedAt).toLocaleDateString()}</td>
                <td className="px-4 py-3 text-gray-600">
                  {bill.periodFrom && bill.periodTo ? `${bill.periodFrom} — ${bill.periodTo}` : 'N/A'}
                </td>
                <td className="px-4 py-3 text-right font-medium text-gray-900">
                  {bill.totalAmount.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 })}
                </td>
                <td className="px-4 py-3 text-right text-gray-600">{bill.transactionCount}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
