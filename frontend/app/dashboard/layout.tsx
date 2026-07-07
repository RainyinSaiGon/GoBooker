import Link from 'next/link';

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex h-screen flex-col md:flex-row md:overflow-hidden bg-zinc-50 dark:bg-zinc-950">
      {/* Sidebar Container */}
      <div className="w-full flex-none md:w-64 border-r border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900">
        <div className="flex h-full flex-col p-4">
          
          {/* Logo / Top Banner */}
          <Link
              href="/"
            >
            <img 
              src="/logo.png" 
              alt="GoBooker Logo" 
              className="h-full w-auto object-contain" 
            />
          </Link>

          {/* Navigation Links and Actions */}
          <div className="flex grow flex-row justify-between space-x-2 md:flex-col md:space-x-0 md:space-y-2">
            {/* Spacer panel matching your design template */}
            <div className="hidden h-auto w-full grow rounded-xl bg-zinc-50 dark:bg-zinc-800/50 md:block"></div>
            
            {/* Action / Sign Out */}
            <form className="w-full">
              <button type="submit" className="flex h-12 w-full items-center justify-center gap-2 rounded-xl bg-zinc-100 hover:bg-red-50 hover:text-red-600 dark:bg-zinc-800 dark:hover:bg-red-950/30 dark:hover:text-red-400 p-3 text-sm font-medium transition-colors md:justify-start">
                <span className="text-sm font-medium">Sign Out</span>
              </button>
            </form>
          </div>
          
        </div>
      </div>

      {/* Dynamic Content Area (Where Dashboard Page loads) */}
      <div className="grow p-6 md:p-12 md:overflow-y-auto flex items-center justify-center">
        {children}
      </div>
    </div>
  );
}