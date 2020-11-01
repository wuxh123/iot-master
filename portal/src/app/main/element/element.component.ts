import {Component, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzModalService, NzTableQueryParams} from 'ng-zorro-antd';
import {Router} from '@angular/router';
import {TabRef} from '../tabs/tabs.component';
import {ElementEditComponent} from "../element-edit/element-edit.component";

@Component({
  selector: 'app-element',
  templateUrl: './element.component.html',
  styleUrls: ['./element.component.scss']
})
export class ElementComponent implements OnInit {

  datum: any[];
  total = 0;
  pageIndex = 1;
  pageSize = 10;
  sortField = null;
  sortOrder = null;
  filters = [];
  keyword = '';
  loading = false;



  constructor(private as: ApiService, private router: Router, private tab: TabRef, private ms: NzModalService) {
    tab.name = '元件管理';
  }

  ngOnInit(): void {
  }

  reload(): void {
    this.pageIndex = 1;
    this.keyword = '';
    this.load();
  }

  load(): void {
    this.loading = true;
    this.as.post('elements', {
      offset: (this.pageIndex - 1) * this.pageSize,
      length: this.pageSize,
      sortKey: this.sortField,
      sortOrder: this.sortOrder,
      filters: this.filters,
      keywords: [
        {key: 'Name', value: this.keyword},
      ]
    }).subscribe(res => {

      this.datum = res.data;
      this.total = res.total;
    }, error => {
      console.log('error', error);
    }, () => {
      this.loading = false;
    });
  }


  edit(id?): void {
    const modal = this.ms.create({
      nzTitle: id ? '编辑元件' : '创建元件',
      nzContent: ElementEditComponent,
      nzFooter: null,
      nzMaskClosable: false,
      // nzViewContainerRef: this.viewContainerRef,
      nzComponentParams: {id},
    });
    // insert/update after close
    modal.afterClose.subscribe(data => {
      if (!data) {
        return;
      }
      if (id) {
        this.datum.forEach((c: any, i, a: any[]) => {
            if (c.id === data.id) {
              a[i] = data;
            }
          }
        );
      } else {
        this.datum.unshift(data);
      }
    });
  }



  detail(c): void {
    this.router.navigate(['/admin/element/' + c.id + '/detail']);
  }

  onTableQuery(params: NzTableQueryParams): void {
    const {pageSize, pageIndex, sort, filter} = params;
    this.pageSize = pageSize;
    this.pageIndex = pageIndex;
    const currentSort = sort.find(item => item.value !== null);
    this.sortField = (currentSort && currentSort.key) || null;
    this.sortOrder = (currentSort && currentSort.value) || null;
    this.filters = filter;
    this.load();
  }

  search(): void {
    this.pageIndex = 1;
    this.load();
  }
}
