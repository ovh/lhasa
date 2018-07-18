import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, ActivatedRoute, Router, Params } from '@angular/router';
import { Observable } from 'rxjs';

@Injectable()
export class RoutingGuard implements CanActivate {

  constructor(
    private route: ActivatedRoute,
    private router: Router
  ) {

  }

  canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return true;
  }
}
