import { useNavigate, Link } from 'react-router-dom';
import { useTransactions } from '../hooks/useTransactions';
import {
  PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer,
  AreaChart, Area, CartesianGrid,
} from 'recharts';

const COLORS = ['#6366f1', '#8b5cf6', '#a78bfa', '#c084fc', '#e879f9', '#f472b6', '#fb7185', '#f87171', '#fb923c', '#fbbf24', '#a3e635', '#34d399'];

const MONTH_NAMES = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

function parseMonthKey(dateStr: string): string {
  const parts = dateStr.split('/');
  if (parts.length < 3) return '';
  const month = parseInt(parts[1], 10);
  const year = parts[2].length === 2 ? `20${parts[2]}` : parts[2];
  return `${year}-${String(month).padStart(2, '0')}`;
}

function formatMonthLabel(key: string): string {
  const [year, month] = key.split('-');
  return `${MONTH_NAMES[parseInt(month, 10) - 1]} ${year}`;
}

export default function Dashboard() {
  const { transactions, categorySummary, totalSpend, dateRange } = useTransactions();
  const navigate = useNavigate();

  if (transactions.length === 0) {
    return (
      <div className="text-center py-20 animate-fade-in-up">
        <div className="text-6xl mb-4">📊</div>
        <h2 className="text-2xl font-bold text-gray-900">No transactions yet</h2>
        <p className="text-gray-500 mt-2">Upload your credit card statements to see spending insights</p>
        <Link to="/upload" className="inline-block mt-6 px-6 py-2.5 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 transition-colors">
          Upload Bills
        </Link>
      </div>
    );
  }

  const pieData = categorySummary.map((c) => ({ name: c.name, value: Math.round(c.total) }));
  const pieTotal = pieData.reduce((s, d) => s + d.value, 0);

  const topSubcategories = categorySummary
    .flatMap((c) => c.subCategories.map((s) => ({ name: `${s.name}`, total: s.total, category: c.name })))
    .sort((a, b) => b.total - a.total)
    .slice(0, 10);

  // Monthly trend data
  const monthlyMap = new Map<string, number>();
  for (const tx of transactions) {
    if (tx.amount <= 0) continue;
    const key = parseMonthKey(tx.date);
    if (key) monthlyMap.set(key, (monthlyMap.get(key) ?? 0) + tx.amount);
  }
  const monthlyTrend = Array.from(monthlyMap.entries())
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([key, total]) => ({ month: formatMonthLabel(key), total: Math.round(total) }));

  const showTrend = monthlyTrend.length >= 2;

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">Dashboard</h2>
      <p className="text-gray-500 mt-1">Your spending overview</p>

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-6">
        {[
          { label: 'Total Spend', value: totalSpend.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }) },
          { label: 'Transactions', value: transactions.length },
          { label: 'Categories', value: categorySummary.length },
          { label: 'Period', value: `${dateRange.from} — ${dateRange.to}`, small: true },
        ].map((card, i) => (
          <div key={i} className="bg-white rounded-xl border border-gray-200 p-5 animate-fade-in-up" style={{ animationDelay: `${i * 80}ms`, animationFillMode: 'backwards' }}>
            <p className="text-sm text-gray-500">{card.label}</p>
            <p className={`${card.small ? 'text-lg' : 'text-2xl'} font-bold text-gray-900 mt-1`}>{card.value}</p>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-8">
        {/* Pie Chart with custom legend */}
        <div className="bg-white rounded-xl border border-gray-200 p-6 animate-fade-in-up" style={{ animationDelay: '320ms', animationFillMode: 'backwards' }}>
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Spend by Category</h3>
          <ResponsiveContainer width="100%" height={260}>
            <PieChart>
              <Pie
                data={pieData}
                cx="50%"
                cy="50%"
                innerRadius={55}
                outerRadius={100}
                paddingAngle={2}
                dataKey="value"
                cursor="pointer"
                onClick={(_, index) => navigate(`/transactions?category=${encodeURIComponent(pieData[index].name)}`)}
              >
                {pieData.map((_, i) => (
                  <Cell key={i} fill={COLORS[i % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip formatter={(value: number) => value.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 })} />
            </PieChart>
          </ResponsiveContainer>
          <div className="mt-3 max-h-40 overflow-y-auto grid grid-cols-2 gap-x-4 gap-y-1.5 text-xs">
            {pieData.map((d, i) => (
              <button
                key={d.name}
                onClick={() => navigate(`/transactions?category=${encodeURIComponent(d.name)}`)}
                className="flex items-center gap-2 text-left hover:bg-gray-50 rounded px-1 py-0.5 transition-colors"
              >
                <span className="w-2.5 h-2.5 rounded-full shrink-0" style={{ backgroundColor: COLORS[i % COLORS.length] }} />
                <span className="text-gray-700 truncate">{d.name}</span>
                <span className="text-gray-400 ml-auto shrink-0">{pieTotal > 0 ? `${Math.round((d.value / pieTotal) * 100)}%` : ''}</span>
              </button>
            ))}
          </div>
        </div>

        {/* Bar Chart */}
        <div className="bg-white rounded-xl border border-gray-200 p-6 animate-fade-in-up" style={{ animationDelay: '400ms', animationFillMode: 'backwards' }}>
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Top Subcategories</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart
              data={topSubcategories}
              layout="vertical"
              margin={{ left: 80 }}
              onClick={(state) => {
                if (state?.activePayload?.[0]?.payload) {
                  const item = state.activePayload[0].payload;
                  navigate(`/transactions?category=${encodeURIComponent(item.category)}&subCategory=${encodeURIComponent(item.name)}`);
                }
              }}
            >
              <XAxis type="number" tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`} />
              <YAxis type="category" dataKey="name" width={80} tick={{ fontSize: 12 }} />
              <Tooltip formatter={(value: number) => value.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 })} />
              <Bar dataKey="total" fill="#6366f1" radius={[0, 4, 4, 0]} cursor="pointer" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Monthly Trend */}
      {showTrend && (
        <div className="mt-8 bg-white rounded-xl border border-gray-200 p-6 animate-fade-in-up" style={{ animationDelay: '480ms', animationFillMode: 'backwards' }}>
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Monthly Spending Trend</h3>
          <ResponsiveContainer width="100%" height={280}>
            <AreaChart data={monthlyTrend} margin={{ left: 20, right: 20, top: 10, bottom: 0 }}>
              <defs>
                <linearGradient id="trendGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#6366f1" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#6366f1" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#f3f4f6" />
              <XAxis dataKey="month" tick={{ fontSize: 12 }} />
              <YAxis tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`} tick={{ fontSize: 12 }} />
              <Tooltip formatter={(value: number) => value.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 })} />
              <Area type="monotone" dataKey="total" stroke="#6366f1" strokeWidth={2} fill="url(#trendGradient)" />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      )}

      {/* Category Breakdown */}
      <div className="mt-8 bg-white rounded-xl border border-gray-200 p-6 animate-fade-in-up" style={{ animationDelay: showTrend ? '560ms' : '480ms', animationFillMode: 'backwards' }}>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Category Breakdown</h3>
        <div className="space-y-3">
          {categorySummary.map((cat, i) => (
            <div
              key={cat.name}
              className="border border-gray-100 rounded-lg p-4 cursor-pointer hover:border-indigo-200 hover:bg-indigo-50/30 transition-colors"
              onClick={() => navigate(`/transactions?category=${encodeURIComponent(cat.name)}`)}
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-3 h-3 rounded-full" style={{ backgroundColor: COLORS[i % COLORS.length] }} />
                  <span className="font-medium text-gray-900">{cat.name}</span>
                  <span className="text-xs text-gray-400">{cat.count} transactions</span>
                </div>
                <span className="font-semibold text-gray-900">{cat.total.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 })}</span>
              </div>
              {cat.subCategories.length > 1 && (
                <div className="mt-2 ml-6 space-y-1">
                  {cat.subCategories.map((sub) => (
                    <div
                      key={sub.name}
                      className="flex justify-between text-sm hover:text-indigo-600 transition-colors"
                      onClick={(e) => {
                        e.stopPropagation();
                        navigate(`/transactions?category=${encodeURIComponent(cat.name)}&subCategory=${encodeURIComponent(sub.name)}`);
                      }}
                    >
                      <span className="text-gray-500">{sub.name} ({sub.count})</span>
                      <span className="text-gray-700">{sub.total.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 })}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
