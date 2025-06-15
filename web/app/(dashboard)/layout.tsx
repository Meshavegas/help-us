import Link from 'next/link';
import { Package2, PanelLeft, Settings } from 'lucide-react';

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from '@/components/ui/breadcrumb';
import { Button } from '@/components/ui/button';
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger
} from '@/components/ui/tooltip';
import { Analytics } from '@vercel/analytics/react';
import { User } from './user';
import { VercelLogo } from '@/components/icons';
import Providers from './providers';
import { SearchInput } from './search';
import { getCurrentUser } from '@/lib/auth';
import { routes as allRoutes, Route as UIRoute } from '@/lib/routes';
import React from 'react';
import DashboardBreadcrumb from '@/components/dashboard/DashboardBreadcrumb';

// Helper: map backend role -> frontend role key used in routes.ts
function mapRole(apiRole?: string): string {
  switch (apiRole) {
    case 'administrator':
      return 'admin';
    case 'enseignant':
      return 'teacher';
    case 'famille':
      return 'family';
    default:
      return apiRole || 'guest';
  }
}

// Recursive filter of UI routes per role
function filterRoutesByRole(list: UIRoute[], role: string): UIRoute[] {
  // admin sees everything
  if (role === 'admin') return list;

  return list
    .map((r) => {
      const children = r.children ? filterRoutesByRole(r.children, role) : [];
      const hasRole = r.roles.includes(role);
      if (!hasRole && children.length === 0) return null;
      return { ...r, children } as UIRoute;
    })
    .filter(Boolean) as UIRoute[];
}

export default async function DashboardLayout({
  children
}: {
  children: React.ReactNode;
}) {
  const user = await getCurrentUser();
  const roleKey = mapRole((user as any)?.role);
  const accessibleRoutes = filterRoutesByRole(allRoutes, roleKey);

  return (
    <Providers>
      <main className="flex min-h-screen w-full flex-col bg-muted/40">
        <DesktopNav routes={accessibleRoutes} />
        <div className="flex flex-col sm:gap-4 sm:py-4 pl-[240px]">
          <header className="sticky top-0 z-30 flex h-14 items-center gap-4 border-b bg-background px-4 sm:static sm:h-auto sm:border-0 sm:bg-transparent sm:px-6">
            <MobileNav routes={accessibleRoutes} />
            <DashboardBreadcrumb />
            <SearchInput />
            <User />
          </header>
          <main className="grid flex-1 items-start gap-2 p-4 sm:px-6 sm:py-0 md:gap-4 bg-muted/40">
            {children}
          </main>
        </div>
        <Analytics />
      </main>
    </Providers>
  );
}

function DesktopNav({ routes }: { routes: UIRoute[] }) {
  return (
    <aside className="fixed inset-y-0 left-0 z-10 hidden w-60 flex-col border-r bg-background sm:flex">
      <div className="flex h-14 items-center px-4 border-b">
        <Link
          href="/dashboard"
          className="flex items-center gap-2 text-lg font-semibold"
        >
          <VercelLogo className="h-5 w-5" /> <span>Dashboard</span>
        </Link>
      </div>
      <nav className="flex-1 overflow-y-auto px-3 py-4">
        {routes.map((route) => {
          if (route.children && route.children.length > 0) {
            return (
              <div key={route.title} className="mb-6">
                <p className="mb-2 px-2 text-xs font-semibold text-muted-foreground uppercase">
                  {route.title}
                </p>
                <div className="space-y-1">
                  {route.children.map((child) => (
                    <Link
                      key={child.href}
                      href={child.href}
                      className="flex items-center gap-2 rounded-md px-2 py-2 text-sm text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                    >
                      {React.createElement(child.icon, { className: 'h-4 w-4' })}
                      <span>{child.title}</span>
                    </Link>
                  ))}
                </div>
              </div>
            );
          }
          // Route sans sous-routes
          return (
            <Link
              key={route.href}
              href={route.href}
              className="flex items-center gap-2 rounded-md px-2 py-2 text-sm text-muted-foreground hover:bg-accent hover:text-accent-foreground mb-1"
            >
              {React.createElement(route.icon, { className: 'h-4 w-4' })}
              <span>{route.title}</span>
            </Link>
          );
        })}
      </nav>
      <div className="border-t p-3">
        <Link
          href="/settings"
          className="flex items-center gap-2 rounded-md px-2 py-2 text-sm text-muted-foreground hover:bg-accent hover:text-accent-foreground"
        >
          <Settings className="h-4 w-4" />
          <span>Settings</span>
        </Link>
      </div>
    </aside>
  );
}

function MobileNav({ routes }: { routes: UIRoute[] }) {
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button size="icon" variant="outline" className="sm:hidden">
          <PanelLeft className="h-5 w-5" />
          <span className="sr-only">Menu</span>
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="sm:max-w-xs">
        <nav className="grid gap-6 text-lg font-medium">
          <Link
            href="/dashboard"
            className="group flex h-10 w-10 shrink-0 items-center justify-center gap-2 rounded-full bg-primary text-lg font-semibold text-primary-foreground md:text-base"
          >
            <Package2 className="h-5 w-5 transition-all group-hover:scale-110" />
            <span className="sr-only">Home</span>
          </Link>
          {routes.map((route) => {
            if (route.children && route.children.length > 0) {
              return (
                <div key={route.title} className="mb-4">
                  <p className="px-2.5 mb-2 text-xs font-semibold uppercase text-muted-foreground">
                    {route.title}
                  </p>
                  {route.children.map((child) => (
                    <Link
                      key={child.href}
                      href={child.href}
                      className="flex items-center gap-4 pl-5 py-2 text-muted-foreground hover:text-foreground"
                    >
                      {React.createElement(child.icon, { className: 'h-5 w-5' })}
                      {child.title}
                    </Link>
                  ))}
                </div>
              );
            }
            return (
              <Link
                key={route.href}
                href={route.href}
                className="flex items-center gap-4 px-2.5 py-2 text-muted-foreground hover:text-foreground"
              >
                {React.createElement(route.icon, { className: 'h-5 w-5' })}
                {route.title}
              </Link>
            );
          })}
        </nav>
      </SheetContent>
    </Sheet>
  );
}