import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';

@Injectable()
export class ConfigurationService {

  public ServerAppCatalogWithUrl: string;
  public ServerCatalogWithApiUrl: string;
  public ServerAppCatalogCompanionWithUrl: string;
  public ServerCatalogCompanionWithApiUrl: string;
  public AuthToken: string;

  private ServerAppCatalogCompanion: string = environment.appcatalogCompanion.baseUrlUi + '/';
  private ApiUrlAppCatalogCompanion: string = environment.appcatalogCompanion.baseUrlUi + '/api/';

  private ServerAppCatalog: string = environment.appcatalog.baseUrlUi + '/';
  private ApiUrlAppCatalog: string = environment.appcatalog.baseUrlUi + '/api/';

  /**
   * constructor
   */
  constructor() {
    this.ServerAppCatalogWithUrl = this.ServerAppCatalog;
    this.ServerCatalogWithApiUrl = this.ApiUrlAppCatalog;
    this.ServerAppCatalogCompanionWithUrl = this.ServerAppCatalogCompanion;
    this.ServerCatalogCompanionWithApiUrl = this.ApiUrlAppCatalogCompanion;
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
