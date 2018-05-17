import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';

@Injectable()
export class ConfigurationService {

  public ApiUrl: string;
  public AuthToken: string;

  /**
   * constructor
   */
  constructor() {
    this.ApiUrl = environment.apiUrl;
  }

  /**
   * fix session token
   * @param AuthToken
   */
  public setAuthToken(AuthToken: string): void {
    this.AuthToken = AuthToken;
  }

  /**
   * get token
   * @param AuthToken
   */
  public getAuthToken(): string {
    return this.AuthToken;
  }

}
