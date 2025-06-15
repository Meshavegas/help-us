'use client';

import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from "../ui/breadcrumb";
import { Link } from "lucide-react";
import React from "react";

export default function DashboardBreadcrumb() {
    const pathname = require('next/navigation').usePathname();
    const segments: string[] = pathname
      .replace(/^\/+|\/+$/g, '') // trim leading/trailing slash
      .split('/')
      .filter(Boolean);
  
    // Afficher "Dashboard" quand on est à la racine
    if (segments.length === 0) {
      return (
        <Breadcrumb className="hidden md:flex">
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbPage>Dashboard</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      );
    }
  
    // Construire les items
    let accumulated = '';
    return (
      <Breadcrumb className="hidden md:flex">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href="/dashboard">Dashboard</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          {segments.map((seg, idx) => {
            accumulated += `/${seg}`;
            const label = seg.replace(/-/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase());
            const isLast = idx === segments.length - 1;
  
            if (isLast) {
              return (
                <React.Fragment key={accumulated}>
                  <BreadcrumbSeparator />
                  <BreadcrumbItem>
                    <BreadcrumbPage>{label}</BreadcrumbPage>
                  </BreadcrumbItem>
                </React.Fragment>
              );
            }
  
            return (
              <React.Fragment key={accumulated}>
                <BreadcrumbSeparator />
                <BreadcrumbItem>
                  <BreadcrumbLink asChild>
                    <Link href={accumulated}>{label}</Link>
                  </BreadcrumbLink>
                </BreadcrumbItem>
              </React.Fragment>
            );
          })}
        </BreadcrumbList>
      </Breadcrumb>
    );
  }
  