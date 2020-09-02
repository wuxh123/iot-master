import {Component, OnInit} from '@angular/core';
import {DomSanitizer, SafeUrl} from "@angular/platform-browser";

@Component({
  selector: 'app-plugin-store',
  templateUrl: './plugin-store.component.html',
  styleUrls: ['./plugin-store.component.scss']
})
export class PluginStoreComponent implements OnInit {

  src: SafeUrl; // = "localhost:4200";

  constructor(private sanitizer: DomSanitizer) {
  }

  ngOnInit(): void {
  }

  test(src): void {
    this.src = this.sanitizer.bypassSecurityTrustResourceUrl(src);
  }

}
